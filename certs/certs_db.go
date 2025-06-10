// Package certs provides functions related to certificates.
// This file, certs_db, is for functions to discover, retrieve, and store cert information in the database.
package certs

import (
	"encoding/json"
	"net/http"
)

// GetServerCert gets the certificate of the requested server, along with the signing chain of the certificate.
func GetServerCert(w http.ResponseWriter, r *http.Request) {
	// TODO: extract the JWT fom the header, check the user and store that information with the certs.
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
	allCerts := Discover(serverName)

	for _, cert := range allCerts {
		certData := ConvertX509ToCertificate(*cert)
		responseCerts = append(responseCerts, &certData)
	}

	// TODO: store all Certificates in responseCerts to DB

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
