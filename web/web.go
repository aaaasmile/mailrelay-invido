package web

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/aaaasmile/mailrelay-invido/conf"
	"github.com/aaaasmile/mailrelay-invido/web/idl"
	"github.com/aaaasmile/mailrelay-invido/web/relay"
	"github.com/aaaasmile/mailrelay-invido/web/srvhandler"
)

func RunService(configfile string) error {

	config, err := conf.ReadConfig(configfile)
	if err != nil {
		return err
	}
	log.Println("Configuration is read")
	serverurl := conf.Current.ServiceURL
	finalServURL := fmt.Sprintf("https://%s", strings.Replace(serverurl, "0.0.0.0", "localhost", 1))
	finalServURL = strings.Replace(finalServURL, "127.0.0.1", "localhost", 1)
	log.Println("Server started with URL ", serverurl)
	log.Println("Try this url: ", finalServURL)

	hwd := srvhandler.SrvHandler{
		Cfg:   config.SecretConfig,
		Debug: config.DebugVerbose,
	}
	srv := &relay.Server{Addr: serverurl, Handler: hwd.MailHandler,
		AuthHandler:  hwd.AuthHandler,
		AuthRequired: true,
		HandlerRcpt:  hwd.RcptHandler,
		TLSListener:  true,
		Appname:      idl.Appname,
		Hostname:     config.SecretConfig.HostName}

	chShutdown := make(chan struct{}, 1)
	go func(chs chan struct{}) {
		var err error
		err = srv.ConfigureTLS("cert/server.crt", "cert/server.key")
		if err == nil {
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Println("Server error. Not listening anymore: ", err)
			chs <- struct{}{}
		}
	}(chShutdown)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt) //We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	log.Println("Enter in server loop")
loop:
	for {
		select {
		case <-sig:
			log.Println("stop because interrupt")
			break loop
		case <-chShutdown:
			log.Println("stop because service shutdown on listening")
			log.Fatal("Force with an error to restart")
			break loop
		}
	}

	log.Println("Bye, service")
	return nil
}
