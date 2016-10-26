package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/awgh/ratnet/api"
	"github.com/awgh/ratnet/nodes/ram"

	"github.com/awgh/bencrypt/ecc"
	//"github.com/awgh/bencrypt/rsa"
	"github.com/awgh/ratnet/policy"
	"github.com/awgh/ratnet/transports/https"
	"github.com/awgh/ratnet/transports/udp"
)

func Test_node_Export_1(t *testing.T) {

	// Content and Routing Key Setup
	routingKey := new(ecc.KeyPair)
	routingKey.GenerateKey()
	contentKey := new(ecc.KeyPair)
	contentKey.GenerateKey()
	node := ram.New(contentKey, routingKey)

	// Profiles
	if err := node.AddProfile("alpha", true); err != nil {
		t.Fatal(err)
	}
	if err := node.AddProfile("beta", false); err != nil {
		t.Fatal(err)
	}

	// Contacts
	key := contentKey.Clone()
	key.GenerateKey()
	if err := node.AddContact("gamma", key.GetPubKey().ToB64()); err != nil {
		t.Fatal(err)
	}
	key.GenerateKey()
	if err := node.AddContact("delta", key.GetPubKey().ToB64()); err != nil {
		t.Fatal(err)
	}

	// Channels
	key.GenerateKey()
	if err := node.AddChannel("epsilon", key.ToB64()); err != nil {
		t.Fatal(err)
	}
	key.GenerateKey()
	if err := node.AddChannel("zeta", key.ToB64()); err != nil {
		t.Fatal(err)
	}

	// Peers
	key.GenerateKey()
	if err := node.AddPeer("eta", true, "eta.url"); err != nil {
		t.Fatal(err)
	}
	key.GenerateKey()
	if err := node.AddPeer("theta", true, "theta.url"); err != nil {
		t.Fatal(err)
	}

	// Router
	node.Router().Patch(api.Patch{From: "one", To: []string{"and", "two"}})
	node.Router().Patch(api.Patch{From: "three", To: []string{"four"}})

	// Policies and Transports
	udpTransport := udp.New(node)
	httpsTransport := https.New("cert.pem", "key.pem", node, true)
	node.SetPolicy(policy.NewPoll(udpTransport, node), policy.NewServer(httpsTransport, ":20001", false))

	// Done, print
	b, err := node.Export()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile("test_config.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(b); err != nil {
		t.Fatal(err)
	}
	defer f.Close()
}

func Test_node_Import_1(t *testing.T) {
	node := ram.New(nil, nil)
	b, err := ioutil.ReadFile("test_config.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := node.Import(b); err != nil {
		t.Fatal(err)
	}
}
