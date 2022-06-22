# todo-app

A todo tracker for learning purposes. https://todo.claudemokbel.com

## develop

### pre-requisites

To develop the code you'll need a Linux or Mac machine and the following software:

* [go 1.18+](https://go.dev/dl/) (backend server)
* [docker](https://docs.docker.com/desktop/) (Postgres database container)
* make (if intending to use the Makefile)

### frontend

#### Available Scripts

In the project directory, you can run:

##### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

##### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

##### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

##### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.


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
$ # you can specify the path to serve frontend assets from my using the --assets flag
$ go run backend/cmd/todo-server/*.go --config /path/to/config.json --assets /path/to/assets/
$ # or use the TODO_ASSETS evironment variable
$ TODO_ASSETS=/path/to/assets go run backend/cmd/todo-server/*.go --config /path/to/config.json
```

To run the backend code with a config file from AWS param store prefix the path to the parameter with 
**awsparamstore://**. The param should be stored as an encrypted string.

#### Creating users

Once the backend is running you can create test users. The examples below use the default example port and API key.

```shell
# create a test user named george
curl -X POST http://localhost:8080/api/users -d '{"name":"george", "password":"password"}'
# fetch all users (requires the server API key)
curl -H "Todo-Api-Key: test" http://localhost:8080/api/users
# login with the newly created user named george and save the cookie in a file name httpcookie
curl -X POST -d '{"name":"george","password":"password"}' http://localhost:8080/api/user/login -c httpcookie
# use the login cookie to read the API key
curl -b httpcookie http://localhost:8080/api/user/key
# read the user info
curl -b httpcookie http://localhost:8080/api/user
# user the user API key to read the user info
curl -H "Authorization: Bearer <apikey>" http://localhost:8080/api/user
# logout and delete the cookie
curl -X DELETE -b httpcookie http://localhost:8080/api/user/logout && rm httpcookie
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
