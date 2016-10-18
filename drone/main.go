package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/awgh/bencrypt/ecc"
	"github.com/awgh/ratnet/nodes/ram"
	"github.com/awgh/ratnet/policy"
	"github.com/awgh/ratnet/transports/https"
)

func main() {

	var dbFile string
	var publicPort, adminPort int

	flag.StringVar(&dbFile, "dbfile", "ratnet.ql", "QL Database File")
	flag.IntVar(&publicPort, "p", 20001, "HTTPS Public Port (*)")
	flag.IntVar(&adminPort, "ap", 20002, "HTTPS Admin Port (localhost)")
	flag.Parse()

	listenPublic := fmt.Sprintf(":%d", publicPort)
	listenAdmin := fmt.Sprintf("localhost:%d", adminPort)

	node := ram.New(new(ecc.KeyPair), new(ecc.KeyPair))

	transportPublic := https.New("cert.pem", "key.pem", node, true)
	transportAdmin := https.New("cert.pem", "key.pem", node, true)

	node.SetPolicy(
		policy.NewServer(transportPublic, listenPublic, false),
		policy.NewServer(transportAdmin, listenAdmin, true))

	log.Println("Public Server starting: ", listenPublic)
	log.Println("Control Server starting: ", listenAdmin)

	node.Start()
}
