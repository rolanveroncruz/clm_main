package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	serverAddress := "www.amazon.com:443"

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", serverAddress, &tls.Config{
		InsecureSkipVerify: true, // WARNING: Do NOT use this in production unless you understand the security implications.
		// For certificate inspection, we might temporarily skip verification to get the raw cert.
		// For a real application, you should handle verification properly.
	})
	if err != nil {
		log.Fatalf("Failed to connect to TLS server: %v", err)
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)

	fmt.Printf("Connected to %s\n", serverAddress)

	peerCertificates := conn.ConnectionState().PeerCertificates

	if len(peerCertificates) == 0 {
		fmt.Println("No peer certificates found.")
		return
	}

	fmt.Println("\n--- Server Certificate Details ---")
	for i, cert := range peerCertificates {
		fmt.Printf("\nCertificate #%d:\n", i+1)
		fmt.Printf("  Subject: %s\n", cert.Subject.CommonName)
		fmt.Printf("  Issuer: %s\n", cert.Issuer.CommonName)
		fmt.Printf("  Serial Number: %s\n", cert.SerialNumber.String())
		fmt.Printf("  Not Before: %s\n", cert.NotBefore.Format(time.RFC3339))
		fmt.Printf("  Not After: %s\n", cert.NotAfter.Format(time.RFC3339))
		fmt.Printf("  DNS Names: %v\n", cert.DNSNames)
		fmt.Printf("  IP Addresses: %v\n", cert.IPAddresses)
		fmt.Printf("  Signature Algorithm: %s\n", cert.SignatureAlgorithm.String())
		fmt.Printf("  Public Key Algorithm: %s\n", cert.PublicKeyAlgorithm.String())
		keySize := getKeySize(cert.PublicKey)
		if keySize > 0 {
			fmt.Printf("  Key Size: %d bits\n", keySize)
		} else {
			fmt.Printf("  Public Key Type: %s (key size unknown)\n", cert.PublicKeyAlgorithm.String())
		}
		// You can access many more fields from the *x509.Certificate struct
		// For example, O, OU, Country, etc.
		fmt.Printf("  Organizations: %v\n", cert.Subject.Organization)
		fmt.Printf("  Organizational Units: %v\n", cert.Subject.OrganizationalUnit)
		fmt.Printf("  Country: %v\n", cert.Subject.Country)

		// To get the raw certificate in PEM format:
		// pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
		// fmt.Printf("  PEM Encoded Certificate:\n%s\n", pemCert)
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
