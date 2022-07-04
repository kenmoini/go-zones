#!/bin/bash

export SERVER_CONFIG_YAML="/etc/go-zones/config.yml"
export GENERATED_DIR="/opt/app-root/generated-conf"

echo -e "\nSTARTING GO ZONES SERVER...\n"
go-zones -mode server -config "${SERVER_CONFIG_YAML}"