package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"ph.certs.com/clm_main/auth"
	"ph.certs.com/clm_main/certs"
	"ph.certs.com/clm_main/middleware"
	"ph.certs.com/clm_main/sql"
	"time"
)

type pkixName struct {
	CommonName         string   `json:"CommonName"`
	Organization       []string `json:"Organization"`
	OrganizationalUnit []string `json:"OrganizationalUnit"`
	Country            []string `json:"Country"`
	Locality           []string `json:"Locality"`
	Province           []string `json:"Province"`
}
type Certificate struct {
	SerialNumber       string    `json:"serial_number"`
	Subject            pkixName  `json:"common_name"`
	Issuer             pkixName  `json:"issuer"`
	NotBefore          time.Time `json:"not_before"`
	NotAfter           time.Time `json:"not_after"`
	SignatureAlgorithm string    `json:"signature_algorithm"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm"`
	PublicKeySize      int       `json:"public_key_size"`
}

func init() {
	sql.InitializeDatabase()
	auth.CreateAdminUser()
}

func main() {
	defer sql.CloseDatabase()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	portStr := ":" + os.Getenv("PORT")

	mux_router := mux.NewRouter()
	mux_router.Handle("/ping", middleware.LoggingMiddleware(http.HandlerFunc(ping)))
	mux_router.Handle("/postping", middleware.LoggingMiddleware(http.HandlerFunc(postPing)))
	mux_router.Handle("/get_server_cert", middleware.LoggingMiddleware(http.HandlerFunc(getServerCert)))
	mux_router.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(auth.Login)))
	println("Listening on port " + portStr + "...")
	err = http.ListenAndServe(portStr, mux_router)
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type StringResponse struct {
		Message string `json:"Message"`
	}

	response := StringResponse{
		Message: "Hello, World!",
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

func postPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type RequestData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type StringResponse struct {
		Message string `json:"Message"`
	}
	var data RequestData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	name := data.Name
	email := data.Email

	response := StringResponse{
		Message: "Hello, " + name + " with email " + email,
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

func getServerCert(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Server string `json:"server"`
	}

	type Response struct {
		Certs []*Certificate `json:"certs"`
	}
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	responseCerts := []*Certificate{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	serverName := requestData.Server
	allCerts := certs.Discover(serverName)

	for _, cert := range allCerts {
		certData := convertToCertificate(*cert)
		responseCerts = append(responseCerts, &certData)
	}
	response := Response{
		Certs: responseCerts,
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

func getKeySize(pubKey interface{}) int {
	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		return pub.N.BitLen() // RSA key size is the bit length of the modulus(N)
	case *ecdsa.PublicKey:
		return pub.Curve.Params().BitSize // ECDSA key size is the bit size of the elliptic curve
	case ed25519.PublicKey:
		return ed25519.PublicKeySize * 8 // Ed25519 key size is fixed(32 bits * 8bits/byte).
	default:
		return 0 // unknown or unsupported public key type
	}
}

func convertToCertificate(from x509.Certificate) Certificate {
	var to Certificate
	to.SerialNumber = from.SerialNumber.String()
	to.Subject.CommonName = from.Subject.CommonName
	to.Subject.Organization = from.Subject.Organization
	to.Subject.OrganizationalUnit = from.Subject.OrganizationalUnit
	to.Subject.Country = from.Subject.Country
	to.Subject.Locality = from.Subject.Locality
	to.Subject.Province = from.Subject.Province
	to.Issuer.CommonName = from.Issuer.CommonName
	to.Issuer.Organization = from.Issuer.Organization
	to.Issuer.OrganizationalUnit = from.Issuer.OrganizationalUnit
	to.Issuer.Country = from.Issuer.Country
	to.Issuer.Locality = from.Issuer.Locality
	to.Issuer.Province = from.Issuer.Province
	to.NotBefore = from.NotBefore
	to.NotAfter = from.NotAfter
	to.SignatureAlgorithm = from.SignatureAlgorithm.String()
	to.PublicKeyAlgorithm = from.PublicKeyAlgorithm.String()
	keySize := getKeySize(from.PublicKey)
	to.PublicKeySize = keySize
	return to
}
