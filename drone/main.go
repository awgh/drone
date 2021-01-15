package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/awgh/ratnet/api"
	"github.com/awgh/ratnet/api/events/defaultlogger"
	"github.com/awgh/ratnet/nodes/ram"

	// Must underscore include any Keypairs, Routers, Policies, or Transports you want compiled in
	_ "github.com/awgh/bencrypt/ecc"
	_ "github.com/awgh/bencrypt/rsa"
	_ "github.com/awgh/ratnet/policy"
	_ "github.com/awgh/ratnet/router"
	_ "github.com/awgh/ratnet/transports/https"
	_ "github.com/awgh/ratnet/transports/tls"
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
	defaultlogger.StartDefaultLogger(node, api.Info)
	node.Start()
	sendInit(node)

	for {
		msg := <-node.Out()
		content := msg.Content.Bytes()
		if err := node.Import(content); err == nil {
			log.Println("Restarting with config:") //, string(content)+"\n")
		} else {
			log.Println("Import failed:", err)
		}
	}
}

func sendInit(node api.Node) {
	peer, _ := node.GetPeer("0")
	if peer != nil {
		log.Println("Attempting Mothership Connection...")
		pub, err := node.CID()
		if err == nil {
			var ipstr string
			for i := 0; i < 5; i++ {
				ip := GetOutboundIP()
				if ip != nil {
					ipstr = ip.String()
					log.Println("Outbound IP: " + ipstr)
					break
				}
				time.Sleep(time.Second)
			}
			log.Println("Sending Init...")
			node.Send("0", []byte(pub.ToB64()+","+ipstr))
		}
	}
}

// GetOutboundIP - Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
