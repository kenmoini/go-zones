#!/bin/bash

export SERVER="${SERVER:="localhost"}"

echo -e "\nPERFORMING RECORD TESTS...\n"

DIG_TEST=$(dig @${SERVER} idm.example.labs +short)
echo "DIG TEST - idm.example.labs: ${DIG_TEST}"
if [ "$DIG_TEST" != "10.12.0.10" ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @${SERVER} priv.example.labs +short)
echo "DIG TEST - priv.example.labs: ${DIG_TEST}"
if [ "$DIG_TEST" != "192.168.0.11" ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @${SERVER} -t AAAA www.example.labs +short)
echo "DIG TEST, AAAA - www.example.labs: ${DIG_TEST}"
if [ "$DIG_TEST" != "fdf4:e2e0:df12:a100::11" ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @${SERVER} -x 10.12.0.10 +short)
echo "DIG TEST - -x 10.12.0.10: ${DIG_TEST}"
if [ "$DIG_TEST" != "idm.example.labs." ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @${SERVER} -x 192.168.0.11 +short)
echo "DIG TEST - -x 192.168.0.11: ${DIG_TEST}"
if [ "$DIG_TEST" != "priv.example.labs." ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @${SERVER} -x fdf4:e2e0:df12:a100::11 +short)
echo "DIG TEST, AAAA - -x fdf4:e2e0:df12:a100::11 : ${DIG_TEST}"
if [ "$DIG_TEST" != "www.example.labs." ]; then
  echo "DIG TEST FAILED"
  exit 1
fi