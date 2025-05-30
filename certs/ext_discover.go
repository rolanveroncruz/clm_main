package certs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"
)

func Discover(server_name string) []*x509.Certificate {
	serverAddress := server_name + ":443"
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

	peerCertificates := conn.ConnectionState().PeerCertificates

	if len(peerCertificates) == 0 {
		fmt.Println("No peer certificates found.")
		return []*x509.Certificate{}
	}
	return peerCertificates

}
