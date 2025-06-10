package certs

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
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

func ConvertX509ToCertificate(from x509.Certificate) Certificate {
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
