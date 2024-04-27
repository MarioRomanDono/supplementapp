# SupplementApp

[![CI action badge](https://github.com/marioromandono/supplementapp/actions/workflows/ci.yaml/badge.svg)](https://github.com/marioromandono/supplementapp/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/marioromandono/supplementapp)](https://goreportcard.com/report/github.com/marioromandono/supplementapp)

**SupplementApp** is an app written in Go for sports supplements management. At the moment, it allows creating, updating, deleting, retrieving by GTIN (Global Trade Identification Number) and listing all the supplements of a database. Right now the only supported database is PostgresSQL, but the project can be expanded to support additional databases if needed.

There are two ways of running this app: as a `net/http` web server, and as AWS Lambda functions (in this case, only finding by GTIN and listing operations are supported). Neither of them are tested in production, so be careful.

This project was developed with the purpose of practising my Go skills, and that's why I'm pretty sure the code can be improved to make it more idiomatic and better. Please feel free to drop any suggestions if you want to :blush:

To test this app, you can run `go test ./...`. This will run every test of the project, including integration and component tests. These tests need Docker to be present in your system, as they use [testcontainers-go](https://golang.testcontainers.org/).

It is also possible to locally run the HTTP server (available in port 8080) by running `make start_server`. In order to start it, Docker and Docker Compose are required to start the database and web server containers, as well as [Goose](https://github.com/pressly/goose) to run the SQL migrations.
