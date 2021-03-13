#!/bin/bash

set -e

echo -e "\nGENERATING ZONES AND CONFIG...\n"
go-zones -mode file -source /etc/go-zones/zones.yml -dir=/opt/app-root/generated-conf

echo -e "\nCOMBINING CONFIGS...\n"
cat /opt/app-root/generated-conf/config/*.internal.forward.conf > /opt/app-root/generated-conf/config/internal-forward-zones.conf || true
cat /opt/app-root/generated-conf/config/*.external.forward.conf > /opt/app-root/generated-conf/config/external-forward-zones.conf || true

cat /opt/app-root/generated-conf/config/*.internal.reverse.conf > /opt/app-root/generated-conf/config/internal-reverse-zones.conf || true
cat /opt/app-root/generated-conf/config/*.external.reverse.conf > /opt/app-root/generated-conf/config/external-reverse-zones.conf || true

echo -e "\nVALIDATING BIND DNS SERVER CONFIGURATION...\n"

if [ ! "$DISABLE_ZONE_CHECKING" == "yes" ]; then /usr/sbin/named-checkconf -z "$NAMEDCONF"; else echo "Checking of zone files is disabled"; fi

echo -e "\nSTARTING BIND DNS SERVER...\n"

export NAMEDCONF=/opt/app-root/vendor/bind/named.conf
export KRB5_KTNAME=/etc/named.keytab
export DISABLE_ZONE_CHECKING=no
export OPTIONS=""


/usr/sbin/named -u named -c ${NAMEDCONF} $OPTIONS -g
