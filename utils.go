package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"time"
)

var HttpClient *fasthttp.Client

func init() {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cloudflareCertificate)

	cert, err := tls.X509KeyPair(cloudflareCertificate, cloudflarePrivateKey)
	if err != nil {
		log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not create a secure tunnel.\n")
	}

	HttpClient = &fasthttp.Client{
		Name:                     "TikTokAutoClaimer",
		NoDefaultUserAgentHeader: false,
		Dial:                     nil,
		DialDualStack:            false,
		TLSConfig: &tls.Config{
			ClientCAs:          caCertPool,
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
			MinVersion:         12,
		},
		MaxConnsPerHost:     1000000,
		MaxIdleConnDuration: 60 * time.Second,
		MaxConnDuration:     60 * time.Second,
		ReadTimeout:         60 * time.Second,
		WriteTimeout:        60 * time.Second,
		MaxConnWaitTimeout:  60 * time.Second,
	}

}

func ReadFile(path string) []string {
	var dataBytes []byte
	var data string
	var err error

	dataBytes, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not open file: %s. Error: %s\n",
			path, err.Error())
	}
	data = string(dataBytes)

	return strings.Split(data, "\r\n")
}

func GetRandomProxy() string {
	if Config.SingleProxy != "" {
		return Config.SingleProxy
	}

	ProxyMutex.Lock()
	defer ProxyMutex.Unlock()

	var randomIndex *big.Int
	var err error

	randomIndex, err = rand.Int(rand.Reader, big.NewInt(int64(len(Proxies)-1)))
	if err != nil {
		log.Fatalf("Could not generate crypto-random number. Error: %s\n", err.Error())
	}

	return Proxies[randomIndex.Int64()]
}

func SendRequest(request *fasthttp.Request) (*fasthttp.Response, error) {
	var response = fasthttp.AcquireResponse()
	var err error

	err = HttpClient.Do(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
