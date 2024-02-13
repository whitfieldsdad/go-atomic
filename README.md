# go-atomic

A Go library for executing tests from the [Atomic Red Team](https://github.com/redcanaryco/atomic-red-team) framework.

⚠️ This project is under active development ⚠️

## Features

- Import tests from the Atomic Red Team framework<sub>1</sub>

<sub>1. Tests can be imported from a directory, from one or more YAML files, or from a tarball.</sub>

## Caveats

- Tests referencing the `atomic-red-team/atomics` folder are not supported.

## Usage

### Command line interface

Throughout this guide, the following commands are equivalent:

```bash
go run main.go
./bin/go-atomic-red-team-darwin-amd64
./bin/go-atomic-red-team-darwin-arm64
./bin/go-atomic-red-team-linux-386
./bin/go-atomic-red-team-linux-amd64
./bin/go-atomic-red-team-linux-arm64
./bin/go-atomic-red-team-windows-amd64
./bin/go-atomic-red-team-windows-arm64
```

#### Importing tasks from the Atomic Red Team framework

[Tests from the Atomic Red Team framework](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics) must be converted into task templates before they can be executed.

To import tests from a directory:

```bash
go run main.go task-templates -w data/content/ generate ~/src/atomic-red-team
```

To import tests from a YAML file in the [format](https://github.com/redcanaryco/atomic-red-team/blob/master/atomic_red_team/atomic_test_template.yaml) used by the Atomic Red Team framework:

```bash
go run main.go task-templates -w data/content/ generate ~/src/atomic-red-team/atomics/T1087.001/T1087.001.yaml
```

To import tests from a tarball file:

```bash
wget -O- -q "https://api.github.com/repos/redcanaryco/atomic-red-team/tarball" > atomic-red-team.tar.gz
go run main.go task-templates -w data/content/ generate atomic-red-team.tar.gz
```
