#!/bin/bash

if [ -z "$1" ]; then
    echo "ERROR: Give API address as param e.g. agency-api.example.com:50051"
    exit 1
fi

rm -rf ./cert
cp -R ./tools/local-env/cert ./cert
rm ./cert/server/server.crt

echo -n | openssl s_client -connect $1 -servername $1 |
    openssl x509 >./cert/server/server.crt
