package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// json results struct for certspotter
type certspotterResult struct {
	Id           string   `json:"id"`
	TbsSha256    string   `json:"tbs_sha256"`
	CertSha256   string   `json:"cert_sha256"`
	DNSNames     []string `json:"dns_names"`
	PubkeySha256 string   `json:"pubkey_sha256"`
	Issuer       struct {
		FriendlyName string   `json:"friendly_name"`
		CaaDomains   []string `json:"caa_domains"`
		PubkeySha256 string   `json:"pubkey_sha256"`
		Name         string   `json:"name"`
	} `json:"issuer"`
	NotBefore time.Time `json:"not_before"`
	NotAfter  time.Time `json:"not_after"`
	Revoked   bool      `json:"revoked"`
}

// json results struct for crt.sh
type crtshResult struct {
	IssuerCAID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int    `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int    `json:"result_count"`
}

// Helpers for storing of unique domains only
var resultMap = make(map[string]bool)
var resultSlice []string

// Function to make a http request and return the body of the response
func fetch(link string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("REQUEST_CREATION_ERROR: %w", err)
	}
	// Attach headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("REQUEST_ERROR: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("READALL_ERROR: %w", err)
	}
	return body, nil
}

// Function to add unique values to the list of results
func add(item string) {
	// regex to match domains
	valid_domain := regexp.MustCompile(`\b(?:https?://)?(?:[a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,5}\b`)
	formatted := valid_domain.FindAllString(item, -1)
	for i := 0; i < len(formatted); i++ {
		if !resultMap[formatted[i]] {
			resultSlice = append(resultSlice, formatted[i])
			resultMap[formatted[i]] = true
		}
	}

}

// Function to recursively fetch results from the certspotter api
// Get an api token from certspotter and set it as the environment variable "CERTSPOTTER"
func recursiveCertspotter(link string) {
	// Variable to keep track of number of iterations when fetching results
	var iter int
	// setup api auth header
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("CERTSPOTTER")),
	}
	// Infinite loop to iterate over every page of results from the api
	for {
		body, err := fetch(link, headers)
		if err != nil {
			log.Fatal(err)
		}
		var response []certspotterResult
		if err := json.Unmarshal(body, &response); err != nil {
			log.Fatal("ERROR DECODING JSON: ", err)
		}
		if len(response) == 0 {
			break
		}
		for _, item := range response {
			for _, dnsName := range item.DNSNames {
				add(dnsName)
				fmt.Println(dnsName)
			}
		}
		// use the id of the last result on the current page as a parameter
		// for the next request
		lastId := response[len(response)-1].Id
		link = fmt.Sprintf("%s&after=%s", link, lastId)
		iter++
	}
	// reset iteration count
	iter = 0
}

func main() {
	// TODO: Add support for input list file of domains
	// get domain from -d parameter
	domain := flag.String("d", "example.com", "A domain to perform a cert transparency lookup on.")
	flag.Parse()
	// Sources
	certspotterUrl := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&expand=dns_names&expand=issuer&expand=issuer.caa_domains", *domain)
	crtshUrl := fmt.Sprintf("https://crt.sh/?q=%s&output=json", *domain)
	// Fetch results
	recursiveCertspotter(certspotterUrl)
	crtsh(crtshUrl)
	// print all results
	for _, v := range resultSlice {
		fmt.Println(v)
	}
}

// Function to fetch results from crt.sh
func crtsh(link string) ([]crtshResult, error) {
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	var items []crtshResult
	// Unmarshal JSON data into the slice
	if err := json.Unmarshal(body, &items); err != nil {
		log.Fatalf("ERROR_UNMARSHALING_JSON: %v", err)
	}
	processCrtshResults(items, err)
	return items, nil
}

// Function to extract domains from crt.sh results
func processCrtshResults(results []crtshResult, err error) {
	if err != nil {
		log.Fatal("Error Fetching results: ", err)
	}
	for _, result := range results {
		if strings.Contains(result.NameValue, "\n") {
			names := strings.Split(result.NameValue, "\n")
			for _, name := range names {
				if !strings.Contains(name, "@") {
					add(name)
				}
			}
		}
	}
}
