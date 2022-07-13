package main

import (
	"Client/proc_manager"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/denisbrodbeck/machineid"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var MachineID = GetHWID()
var AuthenticationClient *http.Client

func init() {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cloudflareCertificate)

	cert, err := tls.X509KeyPair(cloudflareCertificate, cloudflarePrivateKey)
	if err != nil {
		log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not create a secure tunnel.\n")
	}

	AuthenticationClient = &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			TLSClientConfig: &tls.Config{
				ClientCAs:                caCertPool,
				RootCAs:                  caCertPool,
				Certificates:             []tls.Certificate{cert},
				InsecureSkipVerify:       true,
				PreferServerCipherSuites: true,
				MinVersion:               12,
			},
			ForceAttemptHTTP2: false,
		},
		Timeout: 60 * time.Second,
	}
}

func GetHWID() string {
	machineID, err := machineid.ProtectedID("asfa65f1982319f191f651sg61ds9h1d98h49rj984jk9tz%)%/)§(/%)/§$)&/§)$%/§$/Z$//$/$/$/$$Z)§%)(§$($($)§$)%/§$(%/$§&%$!!%§%!§%!§%!§Q4k9z4k9z4j949ad9a484849ef9w4f9w8f4s65v9s8g198g19w1g91d913t13946193619456934166165dfg1d9g9e81f4f9a4dhkjnblkmn0irnm4aisn3nf%/)($/&§()$/&)($§)(§%)§$%$U%)$§Z%)§$Z%)§$HOISDFONSODFN")
	if err != nil {
		log.Fatalf("[WARNING] Could not get HWID quitting.\n")
	}

	return machineID
}

// HWIDCheck is going to get the hardware id every second to prevent tampering even more
func HWIDCheck() {
	for {
		MachineID = GetHWID()
		time.Sleep(3 * time.Second)
	}
}

// ProcessWatchdog keeps track of running procs and checks for any possibly malicious processes
func ProcessWatchdog() {
	var blackList = []string{
		"fiddler", "privoxy", "http debugger",
		"httpdebugger", "wireshark", "mitm",
		"burp", "ngrok", "charles", "ollydbg",
		"ida", "x64dbg",
	}
	var name string
	var procName string

	for {
		procs, _ := proc_manager.Processes()
		for _, proc := range procs {
			procName = strings.ToLower(proc.Executable())

			for _, name = range blackList {
				if strings.Contains(procName, name) {
					proc, err := os.FindProcess(proc.Pid())
					if err == nil {
						_ = proc.Kill()
					}

					log.Fatalln("[WARNING] Possible attack detected! Quitting for your own safety.")
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// LicenseCheckResponse is the way we parse the JSON response of the server
type LicenseCheckResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func AuthenticateLicense() {
	var request = fasthttp.AcquireRequest()
	var response = fasthttp.AcquireResponse()

	defer fasthttp.ReleaseResponse(response)
	defer fasthttp.ReleaseRequest(request)

}

// LicenseCheck checks the license + HWID every few seconds to ensure it is still valid
func LicenseCheck() {
	var response *http.Response
	var jsonResponse = new(LicenseCheckResponse)

	for {

		var authBody = AuthenticationRequest{
			Username: Config.Username,
			License:  Config.License,
			HWID:     MachineID,
		}
		var body, err = json.Marshal(authBody)

		if err != nil {
			if err != nil {
				log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not create auth body: %s\n", err.Error())
			}
		}

		response, err = AuthenticationClient.Post(GetURL("public/login"), "application/json", bytes.NewReader(body))
		if err != nil {
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Sending request: %s\n", err.Error())
		}

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not authenticate connection. Error: %s\n", err.Error())
		}

		err = json.Unmarshal(body, jsonResponse)
		if err != nil {
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not parse the auth server's response. Error: %s\n", err.Error())
		}

		if !jsonResponse.Status {
			if jsonResponse.Message == "" {
				jsonResponse.Message = string(body)
			}
			log.Fatalf("[\x1b[91mERROR\x1b[97m] Could not verify license. Message: %s\n", jsonResponse.Message)
		}

		time.Sleep(15 * time.Second)
	}
}
