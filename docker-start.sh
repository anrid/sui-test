#!/bin/bash

docker build -t sui-local:latest .
docker run --name sui-local --rm -d sui-local:latest
