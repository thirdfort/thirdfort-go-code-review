# Consumer API

Consumer API is an API gateway to multiple services to simplify the API and requests/responses

## What to do?
Review the PR in the interview by talking through the changes. Do not comment on the PR in GitHub.

## Contributing

Please see the CONTRIBUTING.md for PR instructions.

## Running locally

### Running only Consumer API

Note, this requires you to have a running PostgreSQL with values set in either ./config.yaml or in environment variables. See ./internal/config for more information.

```sh
make run
```

To skip authentication between calls (calls to fingerprint) define `NO_AUTH` environment variable or run with
```sh
NO_AUTH=true go run main.go
```

### Running in Docker

This build and runs Consumer API in a container

```sh
make docker:run
```

### Running in Docker with dependencies

This starts PostgreSQL and a mocked version of PA.

```sh
make docker:build
```

To stop, run

```sh
make docker:stop
```

To remove images, run

```sh
make docker:clean
```

To quickly stop, clean and rebuild, run

```sh
make docker:rebuild
```


## Running Tests
### Running Unit Tests
To run tests locally you need to generate the service mocks by running `make generate`.
To remove them you can run `make clean` (also removes binaries etc)

```sh
make generate
make test
```

## Configuration

Rename `./config.yaml-template` to `./config.yaml` and edit it

These can be overridden by environment variables, see `./internal/config/config.go` for details.
All environment variables are prefixed with `CONAPI_` prefix.
