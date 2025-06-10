package certs

import (
	"encoding/json"
	"net/http"
	"ph.certs.com/clm_main/auth"
	"strings"
)

func GetServerCert(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	err := auth.VerifyToken(tokenString)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userEmail, err1 := auth.GetEmailFromToken(tokenString)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	type RequestData struct {
		Server string `json:"server"`
	}

	type Response struct {
		Certs   []*JSONCertificate `json:"certs"`
		Message string             `json:"message"`
	}
	var requestData RequestData
	err3 := json.NewDecoder(r.Body).Decode(&requestData)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errBytes := []byte(err3.Error())
		_, _ = w.Write(errBytes)
		return
	}

	var responseCerts []*JSONCertificate
	responseCerts = make([]*JSONCertificate, 0, 5)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	serverName := requestData.Server
	var response Response
	alreadyPresent := certAlreadyInServer(serverName, userEmail)
	if alreadyPresent {
		response = Response{
			Certs:   nil,
			Message: "Certificate already in server",
		}

	} else {
		allCerts := Discover(serverName)
		for _, cert := range allCerts {
			certData := ConvertX509ToCertificate(*cert)
			responseCerts = append(responseCerts, &certData)
		}

		err2 := insertIntoDB(responseCerts, userEmail)
		if err2 != nil {
			return
		}

		response = Response{
			Certs:   responseCerts,
			Message: "Certificate inserted",
		}

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}
