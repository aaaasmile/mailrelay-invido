# Mail-Relay
Service  usato per mandare mail da remoto passando per un smtp server dall'indirizzo conosciuto
e affidabile.
Il service di questa repository è semplicemente un SMTP mail relay. 
L'interfaccia di collegamento è anch'essa smtp over tls. 
Prerequisito è un valido account email, per esempio gmx che alla fine manda la mail al destinatario.

## Scenario

    ┌─────────────────────┐
    │                     │     ┌──────────────┐
    │                     │     │              │      ┌──────────────┐
    │                     │     │              │      │              │
    │                     │     │            │ │      │              │
    │   IOT device        │     │            └─┼──────┤►             │
    │                     │     │  mail_relay  │      │              │
    │  Crea la mail       │     │              │      │   Gmx service│
    │   via SMTP   ┌──────┼─────┤►  (questo    │      │              │
    │              │      │     │     service) │      │              │
    │                     │     │              │      │              │
    │                     │     │              │      │              │
    └─────────────────────┘     └──────────────┘      │              │
                                                      └──────┬───────┘
                                                             │
                                                             │
                                                             │
                            ┌────────────────┐               │
                            │                │               │
                            │  target@mail.com               │
                            │                │               │
                            │                │               │
                            │ destinatario   ◄───────────────┘
                            │                │
                            │                │
                            └────────────────┘
## Configurazione
Il file di smtp di configurazione si trova in 

    cert/secret-enc.json 
ed è interamente criptato.

## Motivazione
In principio, il dispositivo iot potrebbe mandare direttamente le mail usando il service
gmx senza passare dal relay. Però potrei cambiare la configurazione del relay ed usare
il mio server postfix senza cambiare le credential dei devices che usano questo servizio.

All'inizio ho provato a mandare mails usando gmail da remoto, ma non è stato possibile in modo continuo usando token e auth2.
La ragione principale è che gmail vuole un'autorizzazione manuale dell'uso dell'account di 
posta per mandare mails e il mio token è valido solo per 7 giorni. 
L'utilizzo del service account di g-suite, invece, non invia mails con gmail nella variante free.

Ho considerato la possibilità di settare un server di posta alla Postfix, ma quello che bisogna
configurare e installare per avere un sistema, che in principio invia una mail a settimana,
sembra troppo. 
Dopo aver dato un'occhiata ad un paio di repository tipo
https://github.com/mhale/smtpd e https://github.com/decke/smtprelay dalle quali ho preso gran parte del codice
di questo Mail-Relay e l'ispirazione (vedi "Why another SMTP server?" in smtprelay), ho deciso
di provare un smtp relay per mandare le mie mails saltuarie da dispositivi sparsi in giro.


### Stop del service
Per stoppare il sevice si usa:

    sudo systemctl stop mailrelay-invido

## Deployment su ubuntu direttamente

    cd ~/build/mailrelay-invido
    git pull --all
    ./publish-relay.sh

## Service setup
Ora bisogna abilitare il service:

    sudo systemctl enable mailrelay-invido.service
Ora si fa partire il service (resistente al reboot):

    sudo systemctl start mailrelay-invido
Per vedere i logs si usa:

    sudo journalctl -f -u mailrelay-invido

## Service Config
Questo il conetnuto del file che compare con:

    sudo nano /lib/systemd/system/mailrelay-invido.service
Poi si fa l'enable:

    sudo systemctl enable mailrelay-invido.service
E infine lo start:

    sudo systemctl start mailrelay-invido
Logs sono disponibili con:

    sudo journalctl -f -u mailrelay-invido

Qui segue il contenuto del file mailrelay-invido.service
Nota il Type=idle che è meglio di simple in quanto così 
viene fatto partire quando anche la wlan ha ottenuto l'IP intranet
per consentire l'accesso.

```
[Install]
WantedBy=multi-user.target

[Unit]
Description=mailrelay-invido service
ConditionPathExists=/home/igor/app/go/mailrelay-invido/current/mailrelay-invido.bin
After=network.target

[Service]
Type=idle
User=igor
Group=igor
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/igor/app/go/mailrelay-invido/current/
ExecStart=/home/igor/app/go/mailrelay-invido/current/mailrelay-invido.bin

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/mailrelay-invido
ExecStartPre=/bin/chown igor:igor /var/log/mailrelay-invido
ExecStartPre=/bin/chmod 755 /var/log/mailrelay-invido
StandardOutput=syslog
StandardError=syslog

```

go mod init github.com/aaaasmile/mailrelay-invido


## TLS Server
Per lo sviluppo locale mi serve un server tls. 
keys and certificate:
```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
srv.ListenAndServeTLS("cert/server.crt", "cert/server.key")
```

## Test
Ho testato il relay con gmx e il mio account ventennale di posta senza nessuna difficoltà.

## Secret
Il file con tutte le mail e account è messo nel file secret.json. 
Per funzionare, secret.json deve essere criptato. 
La prima esecuzione genera un file key.pem che viene usato per criptare il secret.  
Si usa -encr alla command line per generare il file.
Poi si fa ripartire mail-relay.
Per aggiornare il server di invido, mi piazzo locale nella dir cert e mando:

    rsync -av *.* <user>@<server>:/home/igor/app/go/mailrelay-invido/current/cert/

