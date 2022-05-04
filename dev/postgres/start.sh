#!/bin/bash

set -e

CONTAINER_NAME="todo-postgres"

POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
DB_USER=dbuser
DB_PASSWORD=dbpassword
DB_NAME=todo

echo "Creating postgres container $CONTAINER_NAME on localhost:5432"

# get the network args for docker run
os=$(uname | tr '[:upper:]' '[:lower:]')
docker_network_args='--network host'
case "$os" in
  linux*)
    os="linux"
    ;;
  darwin*)
    os="darwin"
    docker_network=''
    ;;
  *)
    echo "Unsupported OS: $os"
    exit 2
    ;;
esac

docker run --rm -d $docker_network --name "$CONTAINER_NAME" -p 5432:5432 \
  -v $PWD/dev/postgres/init.sh:/docker-entrypoint-initdb.d/init.sh:ro \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_USER=$POSTGRES_USER \
  -e DB_USER=$DB_USER \
  -e DB_PASSWORD=$DB_PASSWORD \
  -e DB_NAME=$DB_NAME \
  postgres:14-alpine -c TimeZone=UTC