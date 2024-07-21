package main

/*//////////////////////////////////////////////////////////////
                    CERTSPOTTER API LOOKUP
//////////////////////////////////////////////////////////////*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Result []struct {
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

var apiUrl string
var iter = 0
var lastId string

// helpers for storage of unique values only
var resultMap = make(map[string]bool)
var resultSlice = []string{}

func add(item string) {
	if resultMap[item] {
		// Already in the map
		return
	}
	// Add unique value
	resultSlice = append(resultSlice, item)
	resultMap[item] = true
}

func store(resp Result) {
	// The id field of the last issuance is passed to the endpoint in an additional param
	// named "after" with other values remaining the same to work around result pagination
	lastId = resp[len(resp)-1].Id
	after := fmt.Sprintf("&after=%s", lastId)

	if iter != 0 {
		apiUrl += after
	}
	// store unique dns_name output
	for i := 0; i < len(resp); i++ {
		for j := 0; j < len(resp[i].DNSNames); j++ {
			add(resp[i].DNSNames[j])
		}
	}
}

func main() {
	// parse commandline arguments
	domain := flag.String("d", "example.com", "A domain to perform a cert transparency lookup on.")
	flag.Parse()

	apiUrl = fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&expand=dns_names&expand=issuer&expand=issuer.caa_domains", *domain)

	for {
		var body []byte
		func() {
			// send http get request to url
			res, err := http.Get(apiUrl)

			if err != nil {
				log.Fatal("REQUEST_ERROR: ", err)
			}
			defer res.Body.Close()

			body, err = io.ReadAll(res.Body)
			if err != nil {
				log.Fatal("READALL_ERROR:", err)
			}
		}()

		// unmarshall JSON
		var response Result
		if err := json.Unmarshal(body, &response); err != nil {
			log.Fatal("UNMARSHALL_ERROR", err)
		}
		// fmt.Println("RESPONSE LENGTH", len(response))
		// The api returns an empty array when no more results are present.
		// Empty response check enables the break-out from the infinite loop
		if len(response) < 1 {
			break
		}
		iter += 1
		store(response)
		// print output
		for _, v := range resultSlice {
			fmt.Println(v)
		}
	}

}
