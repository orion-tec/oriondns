#!/bin/bash

systemctl stop oriondns

cd backend
go build -o /usr/bin/oriondns ./cmd/oriondns/*.go
cp ./config/production.yaml /etc/oriondns.yaml

systemctl restart oriondns
