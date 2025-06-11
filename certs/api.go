package certs

import (
	"encoding/json"
	"fmt"
	"net/http"
	certsSQL "ph.certs.com/clm_main/certs/sql"
	"ph.certs.com/clm_main/middleware"
	"ph.certs.com/clm_main/sqlite"
)

func GetCertsFromDB(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Certs   []certsSQL.Certificate `json:"certs"`
		Message string                 `json:"message"`
	}

	ctx := r.Context()
	userEmail := middleware.UserFromContext(ctx)
	if userEmail == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	certs, queryErr := sqlite.CertsQueryCental.GetCertificatesFromUserEmail(ctx, userEmail)
	if queryErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := Response{
		Certs:   certs,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("%s", err)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}

}

func GetServerCert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userEmail := middleware.UserFromContext(ctx)
	if userEmail == nil {
		w.WriteHeader(http.StatusBadRequest)
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
	alreadyPresent := certAlreadyInServer(serverName, userEmail.(string))
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

		err2 := insertIntoDB(responseCerts, userEmail.(string))
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
