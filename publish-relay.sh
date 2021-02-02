#!/bin/bash

echo "Builds app"
go build -o mailrelay-invido.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target mailrelay -outdir ~/app/go/mailrelay-invido/zips/
cd ~/app/go/mailrelay-invido/

echo "update the service"
./update-service.sh

echo "Ready to fly"