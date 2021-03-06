#!/bin/bash

echo -e "\nGENERATING ZONES AND CONFIG...\n"
go-zones -mode file -source /etc/go-zones/zones.yml -dir=/opt/app-root/generated-conf

echo -e "\nSTARTING BIND DNS SERVER...\n"

export NAMEDCONF=/etc/named.conf
export KRB5_KTNAME=/etc/named.keytab
export DISABLE_ZONE_CHECKING=no
export OPTIONS=""

if [ ! "$DISABLE_ZONE_CHECKING" == "yes" ]; then /usr/sbin/named-checkconf -z "$NAMEDCONF"; else echo "Checking of zone files is disabled"; fi

/usr/sbin/named -u named -c ${NAMEDCONF} $OPTIONS