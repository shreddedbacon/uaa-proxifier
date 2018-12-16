package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	uaaUrl := os.Getenv("UAA_URL")
	if uaaUrl == "" {
		fatal("You have to set the env var UAA_URL")
	}
	sslCert := os.Getenv("PROXY_SSL_CERT")
	if sslCert == "" {
		fatal("You have to set the env var PROXY_SSL_CERT")
	}
	sslKey := os.Getenv("PROXY_SSL_KEY")
	if uaaUrl == "" {
		fatal("You have to set the env var PROXY_SSL_KEY")
	}
	skipInsecure := false
	if os.Getenv("SKIP_INSECURE") != "" {
		tempInsecure, err := strconv.ParseBool(os.Getenv("SKIP_INSECURE"))
		if err == nil {
			skipInsecure = tempInsecure
		}
	}
	proxy := NewCustomProxy(uaaUrl, skipInsecure)

	http.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServeTLS(":"+port, sslCert, sslKey, nil))
}

func fatalIf(doing string, err error) {
	if err != nil {
		fatal(doing + ": " + err.Error())
	}
}
func fatal(message string) {
	fmt.Fprintln(os.Stdout, message)
	os.Exit(1)
}
