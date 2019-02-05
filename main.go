package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/daragao/slack_faucet/node"
	"github.com/daragao/slack_faucet/server"
)

var buildstamp, githash string

// Config config struct
type config struct {
	PrivateKey string `json:"private-key"`
	NodeURL    string `json:"node-url"`
	ServerPort string `json:"server-port"`
}

func loadConfig(path string) config {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	result := config{}
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func main() {
	c := loadConfig("./config.json")
	serverPort := ":" + c.ServerPort
	nodeURL := c.NodeURL
	privateKey := c.PrivateKey
	fmt.Printf("Server on: %s\n\tBuild Timestamp: %s\n\tGitHash: %s\n", nodeURL, buildstamp, githash)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	client, err := node.New(nodeURL, privateKey[2:])
	if err != nil {
		log.Panic(err)
	}

	_, err = server.New(serverPort, client)
	if err != nil {
		log.Panic(err)
	}
}
