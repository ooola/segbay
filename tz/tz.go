package main

import (
	"fmt"
	"time"
	"strings"
	"math/big"
	"crypto/rand"
)

func getRand(max int64) string {
    r, err := rand.Int(rand.Reader, big.NewInt(max))
    if err == nil {
        return r.String()
    }
    return ""
}

func parseRFC822(s string) {
	t, err := time.Parse(time.RFC822, s) // TODO: add ts offset
	if err != nil {
		fmt.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
		return
	}
	fmt.Printf("t (local): %s\n", t.Local())
}

func parseRFC3339(s string) {
	t, err := time.Parse(time.RFC3339, s) // TODO: add ts offset
	if err != nil {
		fmt.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
		return
	}
	fmt.Printf("%s (RFC3999 local): %s\n", s, t.Local())
}

func parseRFC1123(s string) {
	t, err := time.Parse(time.RFC1123, s) // TODO: add ts offset
	if err != nil {
		fmt.Println("time.Parse error: " + fmt.Sprintf("%+v", err))
		return
	}
	fmt.Printf("%s (RFC1123 local): %s\n", s, t.Local())
}


func getTimeInEbayFormat(plusSeconds uint) string {
	now := time.Now()
	// TODO: account for timezone
	_, offset := now.Zone()
	m := now.Add(time.Duration(plusSeconds)*time.Second)
	m = m.Add(-(time.Duration(offset) * time.Second))
	t := m.Format(time.RFC3339)
	s := strings.Split(t, "-")
	return s[0] + "-" + s[1] + "-" + s[2] + ".000Z"
}


func main() {
	t := time.Now()
	fmt.Printf("now (local): %s\n", t.Local())
	fmt.Println("time in ebay format: " + getTimeInEbayFormat(0))
	parseRFC822("02 Jan 12 15:04 EDT")
	parseRFC822("16 Sep 12 22:04 EDT")
	parseRFC3339("2012-09-18T02:52:14.000Z")
	parseRFC3339("2006-01-02T15:04:05Z-07:00")
	parseRFC3339("2006-01-02T15:04:05Z")
	parseRFC1123("Thu, 05 Jul 2012 23:14:18 GMT")
}
