package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func practice_4() {
	run()
}

const oneMB = 1024 * 1024
const oneGB = 1024 * oneMB
const responseSize = 2 * oneGB

const serverAddr = "localhost:9999"

func startServer() {
	// Mock HTTP server that always returns 2GB of data
	go http.ListenAndServe(serverAddr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-length", fmt.Sprintf("%d", responseSize))

		// 1MB buffer that'll be copied multiple times to the response
		buf := []byte(strings.Repeat("x", oneMB))
		fmt.Println("Starting to send")
		for i := 0; i < responseSize/len(buf); i++ {
			if _, err := w.Write(buf); err != nil {
				log.Fatal("Failed to write to response. Error: ", err.Error())
			}
		}
	}))

	// Some grace period for the server to start
	time.Sleep(100 * time.Millisecond)
}

func run() {
	startServer()

	// HTTP client
	req, err := http.NewRequest("GET", "http://"+serverAddr, nil)
	if err != nil {
		log.Fatal("Error creating HTTP request: ", err.Error())
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making HTTP request: ", err.Error())
	}

	// Read the response header
	fmt.Println("Response: Content-length:", resp.Header.Get("Content-length"))

	bytesRead := 0
	buf := make([]byte, oneMB)

	// Read the response body
	fmt.Println("Starting to read")
	for {
		n, err := resp.Body.Read(buf)
		bytesRead += n

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Error reading HTTP response: ", err.Error())
		}
	}

	fmt.Println("Response: Read", bytesRead, "bytes")
}
