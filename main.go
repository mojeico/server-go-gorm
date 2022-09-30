package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func first(w http.ResponseWriter, req *http.Request) {
	sentry.CaptureMessage("It works!")
	fmt.Println("first")
	fmt.Fprintf(w, "first\n")
}

func main() {

	fmt.Println("START")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://1e5ec765cd894d5187611358861d8bdd@sentry.mojeico.info/2",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,

		//Environment: "dev",
		//Release:     "my-project-name@1.0.0",
		//Debug: true,
		//TracesSampleRate: 1.0,
	})

	if err != nil {
		fmt.Println("ERROR")
		log.Fatalf("sentry.Init: %s", err)
		fmt.Println(err.Error())
		fmt.Println("ERROR")
	}

	defer sentry.Flush(2 * time.Second)

	id := sentry.CaptureMessage("It works! ---- ")

	fmt.Println(id)

	http.HandleFunc("/api/v1/m/hello", hello)
	http.HandleFunc("/api/v1/m/headers", headers)
	http.HandleFunc("/api/v1/m", first)

	http.ListenAndServe(":8080", nil)

	fmt.Println("DONE")

}
