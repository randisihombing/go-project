package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetail(r)
		fmt.Fprintf(w, "Handling incoming orders")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetail(r)
		fmt.Fprintf(w, "Handling users")
	})

	port := 3000

	//load the TLS cert and key
	cert := "cert.pem"
	key := "key.pem"

	//configure TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	//create a custome server
	server := http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	//enable http2
	http2.ConfigureServer(&server, &http2.Server{})

	fmt.Println("Server is Listening to port: ", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Couldn't start server", err)
	}
	//HTTP 1.1 without TLS
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func logRequestDetail(r *http.Request) {
	httpVersion := r.Proto

	fmt.Println("Received request with HTTP version", httpVersion)

	if r.TLS != nil {
		tlsVersion := getTLSVersionName(r.TLS.Version)
		fmt.Println("Received request with TLS version: ", tlsVersion)
	} else {
		fmt.Println("Received request without TLS")
	}

}

func getTLSVersionName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown TLS Version"
	}
}
