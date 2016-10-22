package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/awgh/bencrypt/ecc"
	"github.com/awgh/ratnet/nodes/ram"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "JSON Config File")
	flag.Parse()

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
	}
	node := ram.New(new(ecc.KeyPair), new(ecc.KeyPair))
	err = node.Import(content)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting...\n\n", string(content))
	node.Start()

	//func() {
	for {
		msg := <-node.Out()
		log.Println(msg)
	}
	//}()
}
