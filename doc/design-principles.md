# Design Principles

This document serves to outline any design principles we've actively chosen to develop to to promote long-term clean coding standards and long-term maintenance and evolution of the Consumer API.

## Table of Contents
- [Package Structure](#package-structure)
    - [Web](#web)
    - [Service](#service)
    - [Repository](#repository)
    - [Model](#model)
- [Dependency Injection](#dependency-injection)


## Package Structure
We want to be deliberate in our package structure, promoting basic programming principles like RUM and making it low-cognitive-load for new engineers onboarding or needing to search our service code.

Our package is structured by layer type, and each file is its own domain. We try and split it similar to a multi-layered architecture (example below):


```bash
├── internal
│   ├── web
│   │   └── tasks.go
│   │   └── expectations.go
│   ├── services
│   │   └── tasks.go
│   │   └── expectations.go
│   │   └── mappings.go
│   ├── repositories
│   │   └── database.go
│   │   └── tasks.go
│   │   └── expectations.go
│   ├── models
│   │   └── tasks.go
│   │   └── expectations.go
├── models //public-facing models
│   │   └── tasks.go
│   │   └── expectations.go
```

Note that there is also a `models` package that lives outside of `internal` - this is for your public-facing API models.

### Web
The Web layer is considered the Presentation layer of the service, and is responsible for managing the interface between the server and the application-specific code. For example, it might handle traffic on the route `/transactions/{transaction_id}/tasks`, and manage requests and responses on that endpoint. It'll return application-specific error messages and translate them into server-friendly error messages that the consumer of the API will be able to consume and handle appropriately.

### Service
The Service layer is considered the Application layer of the service, and is responsible for managing any and all business logic for the particular domain. It acts as the interface between the Web layer and the Repository layer, layering in it's business logic in the middle.

### Repository
The Repository layer is considered the Data layer of the service, and is responsible for interacting with the data store. Note I say data store because whilst this is usually the database, in theory this could be any other data store such as GCP Cloud Storage.

### Model
The Model layer is simply any struct definitions that you use in that particular domain. These could be database models (internal structures of how you store your data), or simply structures used for data translation. 


## Dependency Injection
We promote Dependency Injection where we can - this is an incredibly useful concept that adds a small complexity to the code but adds incredible benefits like hugely increased ease of testing, isolation of logic, and promotion of a loosely-coupled codebase.