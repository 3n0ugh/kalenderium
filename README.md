# kalenderium

<p align="right">
  <img src="https://github.com/3n0ugh/kalenderium/blob/main/web/src/assets/kalenderium_logo.png" alt="drawing" width="180" height="160" align="left"/>
</p>

**The Name:**
The Latin word kalenderium, meaning some kind of periodically arranged account book, is applied to medieval manuscript calendars and also to manuals of, for instance, astronomical, astrological, medical, or horticultural ...

**The Project:**
The kalenderium is the sample for a full-stack app capable of adding/deleting/displaying events and user authentication with tokens.

## Architecture

<img width="1214" alt="Screen Shot 2022-03-29 at 20 19 29" src="https://user-images.githubusercontent.com/69458980/160669081-9ae05172-9d77-4fde-82dd-b75664ca7957.png">

The project has three microservices, three different databases, and one frontend application. All coming HTTP requests from
the frontend are handled by the Web API Service. Web API Service talks with the other two services and sends an HTTP response
to the frontend application. The microservices build with go-kit standard library. 

### 1. User Authentication

- The frontend sends an HTTP request to the /v1/signup endpoint. 
- Web API Service handles the request and validates the data.
Then sends a request to Account Service with gRPC. 
- Account Service validates the data and saves user information to the MySQL
database. It then generates a custom authentication token and saves the token to Redis. Then sends it to Web API Service with
a gRPC response.
- Web API Service sends the token as an HTTP response to the frontend. 
- Frontend saves the token to local storage
and redirects the user to the calendar page.

Also, login and logout processes have similar steps.

### 2. Calendar

- The Frontend sends an HTTP request to the /v1/calendar endpoint.
  (The HTTP method can be DELETE, POST, or GET)
- Web API Service handles the request. And to validate user sends an authentication
  token to Account Service with a gRPC request.
- Account service checks the token has a valid type or not. If the token type is valid, Ask the Redis the token exists or not. If the token exists, The Account Service sends a gRPC response to the Web API Service.
- After that, the Web API service sends a gRPC request to Calendar Service.
  (The request change according to the Frontend request's HTTP method)
- Calendar service validates the Event type. If the event type is valid, can do the
  following jobs according to the Frontend request's HTTP Method:
    - Create an event
    - Delete an event
    - Display all events
- Calendar Service talks with the PostgreSQL database and do the required jobs,
  and sends a gRPC response to Web API Service.
- Web API Service sends an HTTP response to Frontend.
- Frontend update the view according to Web API Service response.

## File Structure

```shell
kalenderium
├── cmd
│   ├── account
│   ├── calendar
│   └── web-api
├── internal
│   ├── config
│   ├── context
│   ├── err
│   ├── token
│   └── validator
├── pkg
│   ├── account
│   │   ├── database
│   │   │   └── migrations
│   │   ├── endpoints
│   │   ├── pb
│   │   ├── repository
│   │   ├── store
│   │   └── transport
│   ├── calendar
│   │   ├── database
│   │   │   └── migrations
│   │   ├── endpoints
│   │   ├── pb
│   │   ├── repository
│   │   └── transport
│   └── web-api
│       ├── client
│       ├── endpoints
│       └── transport
└── web
    ├── public
    └── src
        ├── assets
        ├── components
        ├── plugins
        ├── router
        └── views
```
- The Cmd directory includes the service's main files
- The internal directory is Go specific.
  - When the go command sees an import of a package with internal in its path, 
  it verifies that the package doing the import is within the tree rooted at 
  the parent of the internal directory.
  - The project holds the helper modules in the internal directory.
- The pkg directory includes our services.
- The web directory includes our Vue Frontend app.

## Instructions

There are two options:
1) You can run the project without docker.
2) You can run the project with docker.

### 1. Run Project Without Docker

#### Requirements
- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download/)
- [MySQL](https://www.mysql.com/downloads/)
- [make](https://www.gnu.org/software/make/)
- [go-migrate](https://github.com/golang-migrate/migrate)
- [node](https://nodejs.org/en/download/)
- [yarn](https://yarnpkg.com/getting-started/migration)
- [git](https://git-scm.com/downloads)
- Code Editor

#### Steps
- First, clone the repository:
```shell
git clone https://github.com/3n0ugh/kalenderium.git
```
- Move into the kalenderium directory:
```shell
cd kalenderium
```
- Create a PostgreSQL database with a named calendar.
- Create a PostgreSQL superuser named kalenderium.
- Create a MySQL database with a named account.
- Create a MySQL superuser with named kalenderium.
(If you change any name, you need to change configs into api.dev.yaml file.)
- Change the api.dev.yaml configs. Uncomment the 'use without docker' ones
and commented 'the use with docker' ones. 
- Tidy up the go modules:
```shell
go mod tidy --compat=1.17
```
- Create the environment variables for databases:
```shell
export CALENDAR_DB_DSN=postgres://YOUR_DATABASE_USER:YOUR_USER_PASS@localhost/YOUR_DATABASE_NAME?sslmode=disable
export ACCOUNT_DB_DSN=YOUR_DATABASE_USER:YOUR_USER_PASS@/YOUR_DATABASE_NAME
```
- Make database migrations:
```shell
make db/migrate/up/account
make db/migrate/up/calendar
```
- Install the node_modules:
```shell
make vue/install
```
- Run services:
```shell
make local/run/calendar
make local/run/account
make local/run/web-api
```
- Run the Vue app:
```shell
make vue/run
```
- Now, you can use the website from [here](http://localhost:8080).

### 2. Run Project With Docker

Fortunately, we have Container technology. We will see how easy to run the project.

#### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [node](https://nodejs.org/en/download/)
- [yarn](https://yarnpkg.com/getting-started/migration)
- [Docker-Compose](https://docs.docker.com/compose/install/)
- [make](https://www.gnu.org/software/make/)

#### Steps

- First, clone the repository:
```shell
git clone https://github.com/3n0ugh/kalenderium.git
```
- Move into the kalenderium directory:
```shell
cd kalenderium
```
- Install the node_modules:
```shell
make vue/install
```
- Build the containers according to docker-compose.yaml:
```shell
make docker/build
```
- Run the containers:
```shell
make docker/run
```
- After all containers are up, run the Vue app:
```shell
make vue/run
```
Now, you can use the website from [here](http://localhost:8080). <br/>
As you can see, how easy it is :>