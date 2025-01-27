package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gohttp3/client"
	contract "gohttp3/examples/clientServerContract"
	"io"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("-----------------------------HTTP3 Client with no TLS and Quick Config------------------------------")
	http3Client1 := client.CreateHTTP3Client()
	makeCalls(http3Client1)

	// fmt.Println("-------------------------------HTTP3 Client with TLS and Quick Config-------------------------------")
	// tlsConfig := tls.Config{}
	// quicConfig := quic.Config{}
	// http3Client2 := client.CreateHTTP3ConfigClient(&tlsConfig, &quicConfig)
	// makeCalls(http3Client2)
}

func makeCalls(http3Client *http.Client) {
	urls := []string{
		"https://localhost:443/tasks",
		"https://localhost:443/task/4",
		"https://localhost:443/task/create",
		"https://localhost:443/task/5",
		"https://localhost:443/tasks",
	}

	for _, url := range urls {
		var resp *http.Response
		var err error
		if strings.HasSuffix(url, "create") {
			task := contract.Task{
				ID:          "five",
				Description: "Fifth task",
				Completed:   true,
			}
			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(task); err != nil {
				fmt.Println("Error", err)
				return
			}
			payLoad := bytes.NewBuffer(buffer.Bytes())
			resp, err = http3Client.Post(url, "application/json", payLoad)
		} else {
			resp, err = http3Client.Get(url)
		}
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		fmt.Println("URL:", url)
		fmt.Println("HTTP Version:", resp.Proto)
		fmt.Println("Body:", string(body))
	}
}
