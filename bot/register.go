package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)

var (
	latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
)
var cookie = "UUID=d771b0e6-ee5b-b884-adf5-e31e1ed98ee9; my-saved-list=%7B%7D; dtCookie=v_4_srv_7_sn_FBA69B2ED4D4447A44962E674122DC67_perc_100000_ol_0_mul_1_app-3Af9eed1d550ab4737_1; monkey=6e991479-2a8d-46e8-9edd-3044229dacdc; JSESSIONID=F1AE624493BB27053D31628BD5D247B1;"
var awslab = "AWSALB=d90xQ2PrZKEkq7TxPBwemx4qQh4VD4lnMq3c8oKS0/cJLY9dn5kCx8aM9xX2Xxv+R22qhme45g2r9P3CTcXOj4SifeP14O2WwpR1dis5YZtrSLzyVut5rs/a5saa3/WsdMX0DBy4aFZkoFwv5g4uQ30aFkLXkWkJXrV4gahWghGDXq1n10vLubxGjfz9Ig=="

func main() {
	m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)
	purl, _ := url.Parse("http://localhost:8888")
	client := &http.Client{Transport: m.ConfigureTransport(&http.Transport{
		Proxy: http.ProxyURL(purl),
	})}
	getSession((client))
	getCourses("202208", client)
	select {}
}

func getSession(client *http.Client) {
	m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)

	url := "https://app.testudo.umd.edu/services/dropAdd/signOffOtherSessions"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header = http.Header{
		"cache-control":      {"no-cache"},
		"sec-ch-ua":          {`"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`},
		"pragma":             {"no-cache"},
		"if-modified-since":  {"Sat, 01 Dec 2001 00:00:00 GMT"},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
		"sec-ch-ua-platform": {`"macOS"`},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://app.testudo.umd.edu/"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"cookie":             {cookie + ";" + awslab},
		"dnt":                {"1"},
		http.HeaderOrderKey: {
			"cache-control",
			"sec-ch-ua",
			"pragma",
			"if-modified-since",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-platform",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"dnt",
		},
		http.PHeaderOrderKey: m.PseudoHeaderOrder(),
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	//get cookie from headers
	cookiex := resp.Header.Get("Set-Cookie")
	if len(cookiex) == 0 {
		fmt.Println("ERROR GENNING SESSION")
		return
	}
	labtoken := strings.Split(cookiex, ";")[0]
	awslab = labtoken
	fmt.Println("SESSION GENNED")
	// fmt.Println(awslab)
}

func getCourses(termid string, client *http.Client) {
	m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)
	// fmt.Println(cookie)
	// fmt.Println("====================================")
	// fmt.Println(awslab)
	// fmt.Println("====================================")

	// fmt.Println(cookie + ";" + awslab)
	// return
	url := "https://app.testudo.umd.edu/services/dropAdd/regInfo/" + termid
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = http.Header{
		"cache-control":      {"no-cache"},
		"sec-ch-ua":          {`"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`},
		"pragma":             {"no-cache"},
		"if-modified-since":  {"Sat, 01 Dec 2001 00:00:00 GMT"},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
		"sec-ch-ua-platform": {`"macOS"`},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
		"referer":            {"https://app.testudo.umd.edu/"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"cookie":             {cookie + ";" + awslab},
		"dnt":                {"1"},
		http.HeaderOrderKey: {
			"cache-control",
			"sec-ch-ua",
			"pragma",
			"if-modified-since",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-platform",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
			"dnt",
		},
		http.PHeaderOrderKey: m.PseudoHeaderOrder(),
	}

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}
