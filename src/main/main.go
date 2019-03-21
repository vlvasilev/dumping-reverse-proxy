package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(target string) *Prox {
	url, err := url.Parse(target)
	if err != nil {
		log.Println("[ERROR] ", err.Error())
		return nil
	}
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	p.proxy.Transport = &myTransport{}
	p.proxy.ServeHTTP(w, r)
}

type myTransport struct {
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var sb strings.Builder
	sb.WriteString("--------------------------------------Start Communication-------------------------------------------\n")
	dump, err := httputil.DumpRequest(request, true)
	if err == nil {
		sb.WriteString("Request Body:\n")
		sb.WriteString(string(dump))
		sb.WriteString("\n******************************End Request*************************\n\n")
	} else {
		print("\n\nerror in dumb request")
		// copying the response body did not work
		return nil, err
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		print("\n\ncame in error resp here", err)
		return nil, err //Server is not reachable. Server not working
	}

	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		print("\n\nerror in dumb response")
		// copying the response body did not work
		return nil, err
	}

	sb.WriteString("Response Body:\n")
	sb.WriteString(string(body))

	sb.WriteString("\n--------------------------------------End Transaction-------------------------------------------")
	log.Println(sb.String())
	return response, err
}

func main() {
	const (
		defaultPort        = ":9090"
		defaultPortUsage   = "default server port, ':9090'"
		defaultTarget      = "http://127.0.0.1:8080"
		defaultTargetUsage = "default redirect url, 'http://127.0.0.1:8080'"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	redirecturl := flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	fmt.Println("server will run on :", *port)
	fmt.Println("redirecting to :", *redirecturl)

	// proxy
	proxy := NewProxy(*redirecturl)
	if proxy == nil {
		log.Fatal("Can't make proxy.")
	}

	// server redirection
	http.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
