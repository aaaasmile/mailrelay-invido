# mailrelay-invido
Smpt Mail Relay in golang

This is a small Service relay E-Mails via Smtp. Heavily inspired by https://github.com/mhale/smtpd and https://github.com/decke/smtprelay 
(see "Why another SMTP server?" in smtprelay ) adds a service evelope and uses TLS. Smpt credentials are stored in an encrypted file.

I use it for sending E-mails from a Raspberry Pi3 with a relay to a Gmx account.
For example, a small service that send sporadically E-Mails  is the https://github.com/aaaasmile/crawler
