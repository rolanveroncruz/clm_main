package acme

import (
	"fmt"
	"log"
)

type Provider interface {
	Present(domain, token, keyAuth string) error
	CleanUp(domain, token, keyAuth string) error
}

type CustomProviderHttp01 struct {
	Port string
}

// NewCustomProviderHttp01 is the factory for a CustomProviderHttp01
func NewCustomProviderHttp01(port string) (*CustomProviderHttp01, error) {
	return &CustomProviderHttp01{Port: port}, nil
}

func (p *CustomProviderHttp01) Present(domain, token, keyAuth string) error {
	hostAndPort := domain + ":" + p.Port
	jwtToken, err := Login(hostAndPort)
	if err != nil {
		log.Fatal(err)
	}
	putPairErr := PutPair(hostAndPort, jwtToken, token, keyAuth)
	if putPairErr != nil {
		log.Fatal(putPairErr)
	}
	authStr01, getAuthErr := GetAuthString(domain, token)
	if getAuthErr != nil {
		log.Fatal(getAuthErr)
	}
	fmt.Printf("Token: %s\n", token)
	fmt.Printf("Key Auth: %s\n", keyAuth)
	fmt.Printf("Auth String: %s\n", authStr01)
	if keyAuth != authStr01 {
		return fmt.Errorf("KeyAuth does not match Auth String\n")
	}
	fmt.Printf("KeyAuth matches Auth String\n")
	return nil

}

func (p *CustomProviderHttp01) CleanUp(_ string, _ string, _ string) error {
	return nil
}
