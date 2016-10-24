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

	//go func() {
	for {
		msg := <-node.Out()
		nodeNew := ram.New(nil, nil)
		content := msg.Content.Bytes()
		if err := nodeNew.Import(content); err == nil {
			node.Stop()
			node = nodeNew
			log.Println("Restarting...\n\n", string(content)+"\n")
			node.Start()
		} else {
			log.Println("Import failed:\n", string(content)+"\n")
		}
	}
	//}()
}
