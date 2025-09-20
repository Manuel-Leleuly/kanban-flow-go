# Kanban Flow - GO

## Overview

This is RESTful API used for [Kanban Flow - Web](https://www.github.com/Manuel-Leleuly/kanban-flow-web). I'm still learning how to use Go and how to make a RESTful API in general. Therefore, I am open to any suggestions and feedbacks. Have fun.

## How to use

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

### Swagger

Once the project is running, you can try the endpoints either by using the requests from `manual` folder (for example [tickets.http](./test/manual/ticket.http)) ot by accessing the [swagger](http://127.0.0.1:3005/apidocs/index.html)

## Development Requirements

- IDE / Code Editor
- Go (1.24.6 at the time the project was created)
- PostgreSQL (this project uses Postgres from [Neon](https://www.neon.com))

## Environment Variables

| Name             | Optional |
| ---------------- | -------- |
| DB_USER          | no       |
| DB_PASSWORD      | no       |
| DB_DOMAIN        | no       |
| DB_NAME          | no       |
| DB_TEST_NAME     | no       |
| CLIENT_SECRET    | no       |
| DB_USER          | no       |
| ENABLE_DB_LOGGER | yes      |
| LOG_LEVEL        | yes      |
