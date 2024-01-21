#!/bin/bash

RPC="http://127.0.0.1:9000"
FAUCET="http://127.0.0.1:9123/gas"

# Start test validator
docker run -p 9000:9000 -p 9123:9123 --name sui-local -d --rm sui:latest 2>/dev/null

while true
do
    curl -X POST -d '{}' -H 'content-type: application/json' -sSf $FAUCET > /dev/null 2>&1
    res=$?
    echo "Sui test validator status code: ${res}"

    if [ $res = "22" ]; then
        echo -e "server is online"
        break
    else
        echo -e "server is not ready, retrying .."
    fi
    sleep 3
done

# Init Sui client
docker exec sui-local sui client -y envs

# Setup local net
docker exec sui-local sui client new-env --alias local --rpc $RPC
docker exec sui-local sui client switch --env local
docker exec sui-local sui client envs

# Dump Sui client config
echo "client config:"
docker exec sui-local cat /root/.sui/sui_config/client.yaml

# Show all available addresses
myAddr=$(docker exec sui-local bash -c 'sui client addresses --json | jq -r .addresses[0][1]')
echo "My addr: ${myAddr}"

# Call faucet
curl -X POST \
  -d "{\"FixedAmountRequest\":{\"recipient\":\"${myAddr}\"}}" \
  -H 'content-type: application/json' \
  $FAUCET

echo ""

# Check available gas
docker exec sui-local sui client gas
