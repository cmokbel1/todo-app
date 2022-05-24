#!/bin/bash

# check if docker is running
docker stats --no-stream || exit 1

# start container if not running
docker container inspect -f {{.Id}} todo-postgres || bash ./dev/postgres/start.sh

config_file="backend/cmd/todo-server/config.example.json"
if [[ -f "config.json" ]]; then
    config_file="config.json"
fi

go run backend/cmd/todo-server/*.go --config $config_file