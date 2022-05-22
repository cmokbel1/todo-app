# todo-app

A todo tracker for learning purposes.

## develop

### pre-requisites

To develop the code you'll need a Linux or Mac machine and the following software:

* [go 1.18+](https://go.dev/dl/) (backend server)
* [docker](https://docs.docker.com/desktop/) (Postgres database container)
* other things here...

### frontend

Something something requirements, instructions on how to develop.

### backend

Go **1.18+** and **docker** is required to run and develop the backend. 

Go code is formatted using **goimports**, to install it run the following command:

```shell
$ go install golang.org/x/tools/cmd/goimports@latest
```

To run the backend code use the script at **dev/start_backend.sh**.  
This will create a Postgres Docker container and run migrations, creating initial tables and entering any seed data.

```shell
# create a Docker container for Postgres, run the backend and all migrations.
./dev/start_backend.sh
```

To run the backend code with a custom configuration use one of the following commands:

```shell
$ go run backend/cmd/todo-server/*.go --config /path/to/config.json
$ # or use the TODO_CONFIG environment variable
$ TODO_CONFIG=/path/to/config.json go run backend/cmd/todo-server/*.go
```

#### Creating users

Once the backend is running, use the **apiKey** you set in the config file to create test users. The examples below use
the default example port and API key.

```shell
# create a test user named george
curl -X POST -H "Todo-Api-Key: test" http://localhost:8080/api/users -d '{"name":"george", "email":"george@example.com"}'
# fetch all users
curl -H "Todo-Api-Key: test" http://localhost:8080/api/users
```

#### Tests
To run backend tests run one of the following commands:

```shell
# unit tests with coverage
$ go test ./backend/... -cover
# integration tests with coverage (requires running instance of Postgres)
$ go test ./backend/... -cover -tags integration 
```
## deploy

Something about how to deploy the app

## credits

Something about the people that worked on this and their roles.