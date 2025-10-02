# Kanban Flow - GO

## Overview

This is RESTful API used for [Kanban Flow - Web](https://www.github.com/Manuel-Leleuly/kanban-flow-web). I'm still learning how to use Go and how to make a RESTful API in general. Therefore, I am open to any suggestions and feedbacks. Have fun.

## How to use

### Local

If you want to run this project, you can do so by first installing all the dependencies required:

```
go mod download
```

Then, start the development environment by running the following command:

```
go run main.go

// or

go run .
```

You can also use [Air](https://github.com/air-verse/air) to implement live-reload. This project is already initialized therefore you can just run the command:

```
air
```

### Docker

If you want to run this project using docker, you can do so by running the following command:

```
make run
```

### Swagger

Once the project is running, you can try the endpoints either by using the requests from `manual` folder (for example [tickets.http](./test/manual/ticket.http)) ot by accessing the [swagger](http://127.0.0.1:3005/apidocs/index.html)

## Development Requirements

- IDE / Code Editor
- With Docker
  - Docker
  - Docker Compose
  - GNU Make
- Without Docker (versions used at the time the project was created)
  - Go (1.24.6)
  - PostgreSQL (this project uses Postgres from [Neon](https://www.neon.com)). Why Neon you ask? Well I already made the account when I followed programming tutorials on youtube so...

## Environment Variables

| Name              | Optional |
| ----------------- | -------- |
| APP_ENV           | no       |
| DB_USER           | no       |
| DB_PASSWORD       | no       |
| DB_DOMAIN         | no       |
| DB_NAME           | no       |
| DB_TEST_NAME      | yes      |
| PORT              | yes      |
| CLIENT_SECRET     | no       |
| ENABLE_DB_LOGGER  | yes      |
| LOG_LEVEL         | yes      |
| ENABLE_RATE_LIMIT | yes      |
| RATE_LIMIT_RPS    | yes      |
| DB_SSL_MODE       | yes      |
