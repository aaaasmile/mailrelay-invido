package srvhandler

import (
	"bytes"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/aaaasmile/mailrelay-invido/web/relay"
)

func MailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
	// relay to gmail
	remoteHost := "smtp.gmx.net:465"
	hostName := "localhost"
	var auth smtp.Auth
	//host, _, _ := net.SplitHostPort(remoteHost)

	auth = relay.LoginAuth("myemail@gmx.net", "password")

	err := relay.SendMail(
		remoteHost,
		auth,
		from,
		to,
		data,
		hostName,
	)
	if err != nil {
		log.Println("delivery failed", err)
		return err
	}

	log.Printf("%s delivery successful\n", to)

	return nil
}

func RcptHandler(remoteAddr net.Addr, from string, to string) bool {
	//domain = getDomain(to)
	//return domain == "mail.example.com" // could be checked if the sender is on this domain
	log.Println("Rec handler", from, to)
	return true
}

func AuthHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	log.Println("Auth handler")
	return string(username) == "username@example.tld" && string(password) == "password", nil
}
