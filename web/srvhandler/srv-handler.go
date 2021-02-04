package srvhandler

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/aaaasmile/mailrelay-invido/conf"
	"github.com/aaaasmile/mailrelay-invido/web/relay"
)

type SrvHandler struct {
	Cfg   *conf.SecretConfig
	Debug bool
	relay bool
}

func (hw *SrvHandler) MailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Println("Does not appear a mail message")
		return err
	}
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
	remoteHost := hw.Cfg.RemoteSendHost
	hostName := hw.Cfg.HostName
	fmt.Println(string(data[:5000]))
	var auth smtp.Auth
	hw.relay = true
	if hw.relay {
		auth = relay.LoginAuth(hw.Cfg.EMailLogin, hw.Cfg.EmailPassword)
		err := relay.SendMail(
			remoteHost,
			auth,
			from,
			to,
			data,
			hostName,
			hw.Debug,
		)
		if err != nil {
			log.Println("delivery failed", err)
			return err
		}

		log.Printf("%s delivery successful\n", to)
	}

	return nil
}

func (hw *SrvHandler) RcptHandler(remoteAddr net.Addr, from string, to string) bool {
	//return domain == "mail.example.com" // could be checked if the sender is on this domain
	log.Println("Rec handler", from, to)
	return true
}

func (hw *SrvHandler) AuthHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	log.Println("Auth handler")
	return string(username) == hw.Cfg.ServiceUser && string(password) == hw.Cfg.ServicePassword, nil
}
