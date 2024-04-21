package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/biter777/countries"
)

/*
Your task is to programmatically generate a [self-signed] certificate according to the data you receive from the challenge endpoint.

Things you may be asked to include in the certificate:

a specific country as the organization's country
a specific certificate serial number
the domains the certificate should be valid for
specific valid from & to dates
Encode the certificate in DER format with base64 and POST it to the solution endpoint.
*/

type Input struct {
	PrivateKey   string `json:"private_key"`
	RequiredData struct {
		Domain       string `json:"domain"`
		SerialNumber string `json:"serial_number"`
		Country      string `json:"country"`
	} `json:"required_data"`
}

type Output struct {
	Certificate string `json:"certificate"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	decodedPriv, err := base64.StdEncoding.DecodeString(input.PrivateKey)
	p(err)

	priv, err := x509.ParsePKCS1PrivateKey(decodedPriv)
	p(err)
	pub := priv.Public()
	country := input.RequiredData.Country
	countryCode := countries.ByName(country).Alpha2()
	domain := input.RequiredData.Domain
	serialNumber, ok := new(big.Int).SetString(input.RequiredData.SerialNumber, 0)
	if !ok {
		panic(fmt.Errorf("err: unable to create big int from serial_number %s", input.RequiredData.SerialNumber))
	}

	ca := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:    []string{countryCode},
			CommonName: domain,
		},
		DNSNames:              []string{domain},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(5 * time.Minute),
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &ca, &ca, pub, priv)
	p(err)
	derBase64 := base64.StdEncoding.EncodeToString(derBytes)

	output := Output{Certificate: derBase64}
	out, err := json.Marshal(output)
	p(err)

	fmt.Println(string(out))
}
