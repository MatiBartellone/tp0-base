#!/bin/sh

MSG="tp0_echo_healthcheck"
NETWORK="tp0_testing_net"
SERVER_PORT=12345

RESPONSE=$(docker run --rm --network "$NETWORK" busybox sh -c "echo '$MSG' | nc -w 2 server $SERVER_PORT" 2>/dev/null | tr -d '\r\n')

if [ "$RESPONSE" = "$MSG" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
