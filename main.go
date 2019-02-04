package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ResponseJSON struct
type ResponseJSON struct {
	ResponseType string         `json:"response_type,omitempty"`
	Text         string         `json:"text,omitempty"`
	Attachments  []ResponseJSON `json:"attachments,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	//Read the Request Parameter "command"
	command := r.FormValue("command")
	responseURL := r.FormValue("response_url")
	fmt.Println("Hello handler!")
	//Ideally do other checks for tokens/username/etc
	if command == "/faucet" {
		respJSON := ResponseJSON{"ephemeral", "testing the command: " + responseURL, nil}
		payload, err := json.Marshal(respJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))

		go func() {
			time.Sleep(2 * time.Second)
			postResponse(responseURL, ResponseJSON{"ephemeral", "delayed response", nil})
		}()
	} else {
		fmt.Fprint(w, "I do not understand your command.")
	}
}

func postResponse(url string, respJSON ResponseJSON) {
	payload, err := json.Marshal(respJSON)
	if err != nil {
		fmt.Println("ERROR: failed to unmarshal response")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("ERROR: failed to post response")
	}
	fmt.Println(resp)
}

func main() {
	fmt.Println("vim-go")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
