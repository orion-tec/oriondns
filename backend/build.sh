#!/bin/bash

systemctl stop oriondns

go build -o /usr/bin/oriondns ./backend/cmd/oriondns/*.go
cp backend/config/production.yaml /etc/oriondns.yaml

systemctl restart oriondns
