# Ignition

Ignition is a framework for building go applications. It is designed to be
simple, flexible and extensible. 

It provides multiple packages that can be used independently or together to
build applications:

- [application](x/application/README.md): A framework for building
  applications.
- [config](pkg/config/README.md): A layered configuration system that supports
  multiple configuration sources and formats, including environment variables,
  yaml files, and json files.
- [inject](pkg/injector/README.md): A dependency injection system.
- goenv: A package for working with the go environment.
- mongoutil: A set of utilities for working with MongoDB.

## Installation

```bash
export GOPRIVATE=gitlab.com/firestart/*
go get gitlab.com/firestart/ignition
```

## Usage

For usage examples, see the [examples](examples) directory. For more
information, see the individual package READMEs.