package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Create a new HTTP client with OpenTelemetry instrumentation
	client := createClient()

	callService(client, "http://localhost:8882/call/products")
}

func callService(client *http.Client, url string) {
	// Example of making a GET request and add headers
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Failed to call service:", err)
		callService(client, url) // Retry the service call
		return
	}
	defer resp.Body.Close()

	// Process the response (omitted for brevity)
	print("Response status:", resp.StatusCode, "\n")
	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body))
}

func createClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 3 * time.Second,
	}

	// Create a new HTTP client with OpenTelemetry instrumentation
	return &http.Client{
		Transport: otelhttp.NewTransport(netTransport),
		Timeout:   3 * time.Second, // Set a timeout for the client
	}
}
