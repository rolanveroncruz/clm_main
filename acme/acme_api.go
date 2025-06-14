package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"log"
	"net/http"
	"os"
)

// MyUser is the user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func RequestCertificateAcme(w *http.ResponseWriter, req *http.Request) {
	type CertRequest struct {
		Domain string `json:"domain"`
	}
	var certRequest CertRequest
	err := json.NewDecoder(req.Body).Decode(&certRequest)
	if err != nil {
		http.Error(*w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	domain := certRequest.Domain

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: "rolanpaulvc@gmail.com",
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	// config.CADirURL = "http://192.168.99.100:4000/directory"
	///config.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	//err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	http01Provider, err := NewCustomProviderHttp01(":5001")
	if err != nil {
		panic(err)
	}
	err = client.Challenge.SetHTTP01Provider(http01Provider)

	//err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", ""))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)
	privateKeyFile := fmt.Sprintf("certs/%s-private_key.pem", certificates.Domain)
	pkErr := os.WriteFile(privateKeyFile, certificates.PrivateKey, 0644)
	if pkErr != nil {
		fmt.Println(pkErr)
		return
	}
	certFile := fmt.Sprintf("certs/%s-cert.pem", certificates.Domain)
	certErr := os.WriteFile(certFile, certificates.Certificate, 0644)
	if certErr != nil {
		fmt.Println(certErr)
		return
	}
	caCertFile := fmt.Sprintf("certs/%s-cacert.pem", certificates.Domain)
	CaCertErr := os.WriteFile(caCertFile, certificates.IssuerCertificate, 0644)
	if CaCertErr != nil {
		fmt.Println(CaCertErr)
		return
	}
	hostAndPort := domain + http01Provider.Port

	uploadErr := UploadFileToClient(hostAndPort, privateKeyFile)
	if uploadErr != nil {
		return
	}
	uploadErr = UploadFileToClient(hostAndPort, certFile)
	if uploadErr != nil {
		return
	}
	uploadErr = UploadFileToClient(hostAndPort, caCertFile)
	if uploadErr != nil {
		return
	}
}
