package ebay

import (
	"fmt"
	"net"
	"strings"
	"io/ioutil"
	"crypto/tls"
	"net/http"
	"encoding/xml"
)

type getEbayTimeResponse struct {
	XMLName xml.Name `xml:"GeteBayTimeResponse"`
	Timestamp string   `xml:"Timestamp"`
	Ack string   `xml:"Ack"`
	Build string `xml:"Build"`
	Version string `xml:"Version"`
	parsedOk bool
}

func parseResponse(s string) (*getEbayTimeResponse, error){
	r := new(getEbayTimeResponse)
	r.parsedOk = false
	err := xml.Unmarshal([]byte(s), &r)
	if err != nil {
		return nil, err
	}
	r.parsedOk = true
	return r, nil
}

// makes the call to ebay
func ebayCall(xml string, useTLS bool) (string, error) {
	var rv string
	if (useTLS) {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: false},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get("https://example.com")
		fmt.Printf("1st response: %s", resp)
		conn, err := tls.Dial(	"tcp",
		"http://open.api.sandbox.ebay.com/shopping?:443",
		&tls.Config{InsecureSkipVerify: false})
		if err != nil {
			fmt.Println("tls error: ")
			fmt.Println(err)
		}
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		response, err := ioutil.ReadAll(conn)
		fmt.Printf("2nd response: %s", response)
	} else {
		url := "http://open.api.sandbox.ebay.com/shopping?"
		http.Get(url)
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, strings.NewReader(xml))
		// TODO: read these from a configuration file
		req.Header.Add("X-EBAY-API-APP-ID", "segbaye3c-c714-4d11-9312-1eb4b80e36c")
		req.Header.Add("X-EBAY-API-VERSION", "789")
		req.Header.Add("X-EBAY-API-SITE-ID", "0")
		req.Header.Add("X-EBAY-API-CALL-NAME", "GeteBayTime")
		req.Header.Add("X-EBAY-API-REQUEST-ENCODING", "XML")
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		time, err := parseResponse(string(body))
		if err != nil {
			return "", err
		}
		rv = time.Timestamp
	}
	return rv, nil

}

func GetEbayTime() (string, error) {
	xml := `<?xml version="1.0" encoding="utf-8"?>
		<GeteBayTimeRequest xmlns="urn:ebay:apis:eBLBaseComponents">
		</GeteBayTimeRequest>`
	s, err := ebayCall(xml, false)
	if err != nil {
		return "", err
	}
	return s, nil
}

func main() {
	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	response, err := ioutil.ReadAll(conn)
	//fmt.Printf("response: %s", response)

	/*
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	*/
	//r, err := client.Get("https://golang.org/")
	conn2, err := tls.Dial("tcp", "google.com:443", &tls.Config{InsecureSkipVerify: false})
	if err != nil {
		fmt.Println("tls error: ")
		fmt.Println(err)
	}
	fmt.Fprintf(conn2, "GET / HTTP/1.0\r\n\r\n")
	response, err = ioutil.ReadAll(conn2)
	fmt.Printf("2nd response: %s", response)
	GetEbayTime()
}
