package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aaaasmile/mailrelay-invido/crypto"
	"github.com/aaaasmile/mailrelay-invido/web"
	"github.com/aaaasmile/mailrelay-invido/web/idl"
)

const (
	plainsecretname = "secret.json"
	encsecretname   = "secret-enc.json"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var encr = flag.Bool("encr", false, "Encrypt the secret file")

	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}

	priv, pub, err := crypto.GetKeys()
	if err != nil {
		log.Fatal(err)
	}
	if *encr {
		log.Println("Exncrypting the secret file")
		plain, err := ioutil.ReadFile(plainsecretname)
		if err != nil {
			log.Fatalln("Error on readfile ", err)
		}
		enc := crypto.Encrypt(plain, pub)
		log.Printf("File is encrypted to: %v...", enc[:10])

		err = ioutil.WriteFile(encsecretname, enc, 0644)
		if err != nil {
			log.Fatalln("Write file error: ", err)
		}
		log.Println("Enxrypted file is created. Enjoy", encsecretname)
		os.Exit(0)
	} else {
		plain, err := ioutil.ReadFile(encsecretname)
		if err != nil {
			log.Fatalln("Input file error ", err)
		}
		raw, err := crypto.Decrypt(plain, priv)
		if err != nil {
			log.Fatalln("Decrypt file error ", err)
		}
		log.Println("Decrypt the secret success ", encsecretname)
		if err := web.RunService(*configfile, &raw); err != nil {
			panic(err)
		}
	}
}
