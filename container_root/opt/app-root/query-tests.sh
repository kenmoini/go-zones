#!/bin/bash

echo -e "\nPERFORMING RECORD TESTS...\n"

DIG_TEST=$(dig @localhost idm.example.labs +short)
echo "DIG TEST - idm.example.labs: ${DIG_TEST}"
if [ "$DIG_TEST" != "10.12.0.10" ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @localhost priv.example.labs +short)
echo "DIG TEST - priv.example.labs: ${DIG_TEST}"
if [ "$DIG_TEST" != "192.168.0.11" ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @localhost -x 10.12.0.10 +short)
echo "DIG TEST - -x 10.12.0.10: ${DIG_TEST}"
if [ "$DIG_TEST" != "idm.example.labs." ]; then
  echo "DIG TEST FAILED"
  exit 1
fi

DIG_TEST=$(dig @localhost -x 192.168.0.11 +short)
echo "DIG TEST - -x 192.168.0.11: ${DIG_TEST}"
if [ "$DIG_TEST" != "priv.example.labs." ]; then
  echo "DIG TEST FAILED"
  exit 1
fi