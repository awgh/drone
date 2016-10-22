package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/awgh/ratnet/nodes/ram"

	// Must underscore include any Keypairs, Routers, Policies, or Transports you want compiled in
	_ "github.com/awgh/bencrypt/ecc"
	_ "github.com/awgh/ratnet/policy"
	_ "github.com/awgh/ratnet/router"
	_ "github.com/awgh/ratnet/transports/https"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "JSON Config File")
	flag.Parse()

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
	}
	node := ram.New(nil, nil)
	err = node.Import(content)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting...\n\n", string(content))
	node.Start()

	//go func() {
	for {
		msg := <-node.Out()
		log.Println(msg)
	}
	//}()
}
