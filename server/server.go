package server

import (
	"bytes"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/daragao/slack_faucet/node"
)

// Server base type
type Server struct {
	node *node.EthInstance
}

// ResponseJSON struct
type ResponseJSON struct {
	ResponseType string         `json:"response_type,omitempty"`
	Text         string         `json:"text,omitempty"`
	Attachments  []ResponseJSON `json:"attachments,omitempty"`
}

// New server
func New(port string, node *node.EthInstance) (*Server, error) {
	s := new(Server)
	s.node = node
	http.HandleFunc("/", s.handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return s, nil
}

func sendResponse(w http.ResponseWriter, msgType, msg string) {
	respJSON := ResponseJSON{msgType, msg, nil}
	payload, err := json.Marshal(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("sendResponse error: %s\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

// Handler for all the command requests
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	command := r.FormValue("command")
	responseURL := r.FormValue("response_url")
	text := strings.TrimSpace(r.FormValue("text"))

	if command == "/faucet" {
		if !node.IsHexAddress(text) {
			errMsg := "text needs to be a valid hex address"
			sendResponse(w, "ephemeral", errMsg)
			log.Printf("Faucet failed: %s: %s\n", errMsg, text)
			return
		}

		sendResponse(w, "ephemeral", "Creating Tx to fund account "+text)
		log.Println("Faucet fund:", text)

		go func(address string) {
			time.Sleep(2 * time.Second)
			txHash, err := s.node.Faucet(text, big.NewInt(1000000000000000000)) // 1 eth
			if err != nil {
				log.Printf("Faucet fund failed: %s %s\n", text, err)
				postResponse(responseURL, ResponseJSON{"ephemeral", "failed to fund: " + address, nil})
				return
			}
			postResponse(responseURL, ResponseJSON{"ephemeral", "Submited Tx(" + txHash + ") to fund account" + address, nil})
		}(text)
	} else {
		sendResponse(w, "ephemeral", "\""+command+"\" command not recognized!")
	}
}

func postResponse(url string, respJSON ResponseJSON) {
	payload, err := json.Marshal(respJSON)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Failed to post response: ", err)
	}

	if resp.StatusCode != 200 {
		log.Println("Something went wrong with the delayed post: ", resp.Status)
	}
}
