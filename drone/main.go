package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/awgh/ratnet/nodes/ram"

	// Must underscore include any Keypairs, Routers, Policies, or Transports you want compiled in
	_ "github.com/awgh/bencrypt/ecc"
	_ "github.com/awgh/bencrypt/rsa"
	_ "github.com/awgh/ratnet/policy"
	_ "github.com/awgh/ratnet/router"
	_ "github.com/awgh/ratnet/transports/https"
	_ "github.com/awgh/ratnet/transports/udp"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "JSON Config File")
	flag.Parse()

	// Load initial configuration from file
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	node := ram.New(nil, nil)
	if err := node.Import(content); err != nil {
		log.Fatal(err)
	}
	log.Println("Starting...\n\n", string(content)+"\n")
	node.Start()

	for {
		msg := <-node.Out()
		content := msg.Content.Bytes()
		if err := node.Import(content); err == nil {
			log.Println("Restarting with config:\n\n") //, string(content)+"\n")
		} else {
			log.Println("Import failed:", err)
		}
	}
}
