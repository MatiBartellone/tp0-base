#!/bin/bash

set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Uso: $0 <archivo_salida> <cantidad_clientes>"
  exit 1
fi

OUTPUT_FILE="$1"
CLIENT_COUNT="$2"

if ! [[ "$CLIENT_COUNT" =~ ^[0-9]+$ ]]; then
  echo "Error: la cantidad de clientes debe ser un entero mayor o igual a 0"
  exit 1
fi

cat > "$OUTPUT_FILE" <<EOF
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - SERVER_EXPECTED_AGENCIES=$CLIENT_COUNT
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net
EOF

for ((i=1; i<=CLIENT_COUNT; i++)); do
  cat >> "$OUTPUT_FILE" <<EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data:/.data
    networks:
      - testing_net
    depends_on:
      - server
EOF
done

cat >> "$OUTPUT_FILE" <<EOF

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
EOF
