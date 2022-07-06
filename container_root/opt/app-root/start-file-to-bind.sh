#!/bin/bash

set -e

export NAMEDCONF=/opt/app-root/vendor/bind/named.conf
export KRB5_KTNAME=/etc/named.keytab
export DISABLE_ZONE_CHECKING=no
export OPTIONS=""
export SERVER_CONFIG_YAML="/etc/go-zones/server.yml"
export GENERATED_DIR="/opt/app-root/generated-conf"
export HALT_STARTUP="false"

## Useful for runtime container debugging
while [ "$HALT_STARTUP" == "true" ]; do
  sleep 3600
done

echo -e "\nGENERATING ZONES AND CONFIG...\n"
go-zones -mode file -source "${SERVER_CONFIG_YAML}" -dir "${GENERATED_DIR}"

#echo -e "\nCOMBINING CONFIGS...\n"
#cat /opt/app-root/generated-conf/config/*.internal.forward.conf > /opt/app-root/generated-conf/config/internal-forward-zones.conf || true
#cat /opt/app-root/generated-conf/config/*.external.forward.conf > /opt/app-root/generated-conf/config/external-forward-zones.conf || true
#
#cat /opt/app-root/generated-conf/config/*.internal.reverse.conf > /opt/app-root/generated-conf/config/internal-reverse-zones.conf || true
#cat /opt/app-root/generated-conf/config/*.external.reverse.conf > /opt/app-root/generated-conf/config/external-reverse-zones.conf || true

echo -e "\nVALIDATING BIND DNS SERVER CONFIGURATION...\n"

if [ ! "$DISABLE_ZONE_CHECKING" == "yes" ]; then /usr/sbin/named-checkconf -z "$NAMEDCONF"; else echo "Checking of zone files is disabled"; fi

echo -e "\nSTARTING BIND DNS SERVER...\n"

/usr/sbin/named -u named -c ${NAMEDCONF} $OPTIONS -g