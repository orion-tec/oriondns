#!/bin/bash

systemctl stop oriondns

cd backend
go build -o /usr/bin/oriondns ./cmd/dnsserver/*.go
cp ./config/staging.yaml /etc/oriondns.yaml

systemctl restart oriondns
