package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	runtimeAPI := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	if runtimeAPI == "" {
		panic("AWS_LAMBDA_RUNTIME_API is not set")
	}
	target := fmt.Sprintf("http://%s/2015-03-31/functions/function/invocations", runtimeAPI)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received: %s %s\n", r.Method, r.URL.Path)
		headers := map[string]string{}
		for name, values := range r.Header {
			for _, value := range values {
				headers[name] = value
			}
		}

		queryStringParameters := map[string]string{}
		for name, values := range r.URL.Query() {
			for _, value := range values {
				queryStringParameters[name] = value
			}
		}

		sendEvent := events.APIGatewayProxyRequest{
			HTTPMethod:            r.Method,
			Path:                  r.URL.Path,
			Headers:               headers,
			QueryStringParameters: queryStringParameters,
			RequestContext:        events.APIGatewayProxyRequestContext{},
		}
		sendEventJSON, err := json.Marshal(sendEvent)
		if err != nil {
			panic(err)
		}

		fmt.Printf("-> Forwarding: POST %s - %s\n", target, string(sendEventJSON))

		forwardResponse, err := http.Post(target, "application/json", bytes.NewBuffer(sendEventJSON))
		if err != nil {
			panic(err)
		}
		responseEvent := events.APIGatewayProxyResponse{}
		defer forwardResponse.Body.Close()
		responseEventJSON, err := io.ReadAll(forwardResponse.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(responseEventJSON, &responseEvent)
		if err != nil {
			panic(err)
		}

		fmt.Printf("<- Forwarding Response: %d - %s\n", responseEvent.StatusCode, responseEventJSON)
		fmt.Printf("<- Sending: %d - %s\n", responseEvent.StatusCode, responseEvent.Body)
		w.WriteHeader(responseEvent.StatusCode)
		for name, value := range responseEvent.Headers {
			w.Header().Set(name, value)
		}
		w.Write([]byte(responseEvent.Body))
	})

	fmt.Printf("Listening on %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic(err)
	}
}
