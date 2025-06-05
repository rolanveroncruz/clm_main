package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	authSQL "ph.certs.com/clm_main/auth/sql"
	"ph.certs.com/clm_main/sqlite"
	"time"
)

var SecretKey = []byte("mySecret-Key")

func init() {
	sqlite.InitializeDatabase()
}

/*
CreateAdminUser creates the first admin user, rolanvc@gmail.com
with password <PASSWORD>. This should do this only if the user table is empty or if rolanvc@gmail.com does not yet exist
in the table of users.
*/
func CreateAdminUser() {
	var err error
	adminUserExists, err := checkIfAdminUserExists()
	if err != nil {
		panic(err)
	}
	if adminUserExists {
		println("admin user exists, skipping creation.")
		return
	}
	ctx := context.Background()
	passwordRaw := "<PASSWORD>"
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), 14)
	hashedPassword := string(hash)
	newUUID := uuid.New()
	uUUID := newUUID.String()
	adminUser, err := sqlite.QueryCental.CreateUsers(ctx, authSQL.CreateUsersParams{
		Name:     sql.NullString{String: "RolanVC", Valid: true},
		Email:    sql.NullString{String: "rolanvc@gmail.com", Valid: true},
		Password: sql.NullString{String: hashedPassword, Valid: true},
		Uuid:     sql.NullString{String: uUUID, Valid: true},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("created: %s\n", Scan(adminUser.Name))
}

/*
checkIfAdminUserExists() checks if the "rolanvc@gmail.com" user is already in the database.
*/
func checkIfAdminUserExists() (bool, error) {
	var err error
	ctx := context.Background()
	rowCount, err := sqlite.QueryCental.GetUserCount(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false, errors.New(err.Error())
	}
	if rowCount == 0 {
		return false, nil
	}
	adminUser, err := sqlite.QueryCental.GetUser(ctx, sql.NullString{String: "rolanvc@gmail.com", Valid: true})
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false, errors.New(err.Error())
	}
	if adminUser.Email.Valid {
		return true, nil
	} else {
		return false, errors.New("email is invalid")
	}
}

/*
Scan takes a sql.NullString and extracts the string inside.
*/
func Scan(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "<null>"
}

/*
Login is a handler for logins. (e.g. /login).
*/
func Login(w http.ResponseWriter, r *http.Request) {
	catchAllErrorString := "Invalid Email or Password"
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type LoginResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}

	var requestData LoginRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	email := requestData.Email
	User, err := sqlite.QueryCental.GetUser(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		http.Error(w, catchAllErrorString, http.StatusBadRequest)
		return
	}
	if !User.Email.Valid {
		http.Error(w, catchAllErrorString, http.StatusBadRequest)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.Password.String), []byte(requestData.Password))
	if err == nil {
		tokenString, err := createToken(email, Scan(User.Name))
		response := LoginResponse{
			Name:  User.Name.String,
			Email: User.Email.String,
			Token: tokenString,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	} else {
		http.Error(w, catchAllErrorString, http.StatusBadRequest)
		return
	}
}

func createToken(userEmail string, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userEmail,
		"name":  userName,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
