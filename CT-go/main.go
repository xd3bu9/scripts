package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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

var resultMap = make(map[string]bool)
var resultSlice []string

func fetch(link string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("REQUEST_CREATION_ERROR: %w", err)
	}
	// Add headers
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

func add(item string) {
	if !resultMap[item] {
		resultSlice = append(resultSlice, item)
		resultMap[item] = true
	}
}

func recursiveCertspotter(link string) {
	var iter int
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("CERTSPOTTER")),
	}
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
		lastId := response[len(response)-1].Id
		link = fmt.Sprintf("%s&after=%s", link, lastId)
		iter++
	}
	iter = 0
}

func main() {
	domain := flag.String("d", "example.com", "A domain to perform a cert transparency lookup on.")
	flag.Parse()
	certspotterUrl := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&expand=dns_names&expand=issuer&expand=issuer.caa_domains", *domain)
	crtshUrl := fmt.Sprintf("https://crt.sh/?q=%s&output=json", *domain)

	recursiveCertspotter(certspotterUrl)
	crtsh(crtshUrl)

	for _, v := range resultSlice {
		fmt.Println(v)
	}
}

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
