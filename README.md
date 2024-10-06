# Golinks

[![Go Report Card](https://goreportcard.com/badge/github.com/reimirno/golinks)](https://goreportcard.com/report/github.com/reimirno/golinks)
[![codecov](https://codecov.io/github/Reimirno/golinks/branch/main/graph/badge.svg?token=WFR37HA0LH)](https://codecov.io/github/Reimirno/golinks)
[![CI](https://github.com/Reimirno/golinks/actions/workflows/makefile-ci.yml/badge.svg)](https://github.com/Reimirno/golinks/actions/workflows/makefile-ci.yml)
[![Build](https://github.com/Reimirno/golinks/actions/workflows/makefile-build.yml/badge.svg)](https://github.com/Reimirno/golinks/actions/workflows/makefile-build.yml)

Golinks is a keyword-to-URL mapping service. It allows the user to enter `go/<keyword>` in their browser, and is redirected to the corresponding URL.

This services can be used for personal purposes (e.g. `go/me` opens your website) or for internal link shortening in organizations (e.g. `go/kpi` opens the company's KPIs dashboard).

## Components

This repository consists of a few services:

Server-side services:
- `redirector` - a simple HTTP server that redirects keyword requests to the corresponding URL.
- `crud` - a gRPC service that provides CRUD operations for the keyword-to-URL mappings.
- `crud_http` - an HTTP service that provides CRUD operations for the keyword-to-URL mappings.
- `web` - a web app that allows easier management of the mappings and the server. (WIP)
- `browser/chrome` - a Chrome extension that allows the user to configure the Golinks server URL.

The mapping can be stored in various ways including local files or in database. This behavior is configurable via a configuration file.

## Getting started

You can run both services with all default configurations by:

```bash
go run .
```

Then, go to your browser and try open `localhost:8080/gh`. It should redirect you to `https://github.com`.

## Chrome Extension
If you would like to just type in `go/gh` (instead of `localhost:8080/go/gh`) and make it work, you need to add `127.0.0.1 go` to your host file, which usually is `/etc/hosts` on Linux/MacOS and `C:\Windows\System32\drivers\etc\hosts` on Windows. 

If you are using Chrome, you can also install the extension that is included in `browser/chrome/src` directory to your Chrome browser. Basically, go to `chrome://extensions`, enable "Developer mode" (using top-right corner switch), click "Load unpacked" and select the `browser/chrome/src` directory.

Then, in extension configuration page, you can specify the server URL where you are hosting the `redirector` service, which by default is `http://localhost:8080`.

## Configuration File

You can specify a configuration file by running:

```bash
go run . -config <path_to_config_file>
```

See `files/config.yaml` for schema.

## Mappers

Mappers are key-value stores that map keywords to URLs.

You can specify the mapper in the configuration file. The redirector services supports 4 types of mappers:

| type   | description                                           | configuration      | singleton | readonly |
| ------ | ----------------------------------------------------- | ------------------ | --------- | -------- |
| memory | stores mapping in memory                              | pairs              | true      | true     |
| file   | stores mapping in a local file                        | path, syncInterval | false     | true     |
| bolt   | stores mapping in bolt.db (local file-based kv store) | path, timeout      | true      | false    |
| sql    | stores mapping in a SQL database                      | driver, dsn        | true      | true     |

`readonly` mappers does not support put or delete operations.
`singleton` mappers can only exist once in the system. You can specify one single such mapper in the configuration file.

You can only designate one mapper as `persistor` in the configuration file, and that mapper has to be a non-readonly mapper. 

If the `persistor` field is not specified, then the entire system would be readonly.

## Conflict resolution

If there are multiple mappers configured, CRUD operations would be resolved by the following rules:
- GET: any query would run through the list of mappers and the first successful match would be returned.
- PUT: 
    - insert: the mapper designated as `persistor` would be used.
    - update: the key-value pair would be updated in all mappers that support the GET operation.
- DELETE: the key-value pair would be deleted from the mapper that returns a match by the GET operation rule.

## Sanitization

See code comments in `pkg/sanitizer` for details.

## CRUD gRPC service

The CRUD operations are exposed as a gRPC service. You can use the `grpcurl` tool to interact with the service.

```bash
grpcurl -plaintext -d '{"path": "gh"}' localhost:8081 pb.Golinks/GetUrl
```

## CRUD HTTP service

The CRUD operations are also exposed as a HTTP service. You can use the `curl` tool to interact with the service.

```bash
curl -v http://localhost:8082/go
curl -v http://localhost:8082/go/gh
curl -X PUT -v localhost:8082/go \
     -H "Content-Type: application/json" \
     -d '{"path":"prom","url":"https://prometheus.io"}'
curl -X DELETE -v http://localhost:8082/go/prom
```

This is intended to be interacted with by a CLI tool.

## Web interface

WIP. Should provide:
- basic CRUD operations for user to manage the URL mappings.
- listing all the existing keywords, the list should be searchable and sortable (esp by use count).
- route go/d/* to the web interface. (need to change chrome extension code)
- fuzzy search for keywords (need to implement in backend too)

## Developing

- `redirector`,`crud`,`crud_http`: See `Makefile` for commands to run tests, build and clean the server project.
- `web`: Go to `web`. `npm i` to install dependencies. `npm run dev` to start the development server.
- `browser/chrome`: Go to `browser/chrome`. `npm i` to install dependencies. `npm test` to run the tests. See chrome doc to load the extension in browser.

## Future work
- Web UI/CLI for easier management of the mappings.
- Fuzzy search for keywords
- Deployment scheme
    - containerize and use Kubernetes, Terraform for deployment. Will be more necessary if we want to scale/use stuff like envoy (for grpc-web proxying for example) or connecting to logging/monitoring services.
- Authentication/Authorization
    - integrate with org-specific auth services for access control (who gets to modify/view the mappings etc).
- Monitoring and logging
    - the logs are now going into stdout only.
