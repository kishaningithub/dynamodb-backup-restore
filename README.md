# dynamodb-backup-restore

[![Build Status](https://travis-ci.org/kishaningithub/dynamodb-backup-restore.svg?branch=master)](https://travis-ci.org/kishaningithub/dynamodb-backup-restore)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/kishaningithub/dynamodb-backup-restore)](https://goreportcard.com/report/github.com/kishaningithub/dynamodb-backup-restore)
[![Downloads](https://img.shields.io/github/downloads/kishaningithub/dynamodb-backup-restore/latest/total.svg)](https://github.com/kishaningithub/dynamodb-backup-restore/releases)
[![Latest release](https://img.shields.io/github/release/kishaningithub/dynamodb-backup-restore.svg)](https://github.com/kishaningithub/dynamodb-backup-restore/releases)

A no sweat backup and restore tool for dynamodb

## Table of Contents

- [dynamodb-backup-restore](#dynamodb-backup-restore)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Using Homebrew](#using-homebrew)
    - [Using Binary](#using-binary)
  - [Example](#example)
  - [Usage](#usage)
  - [Maintainers](#maintainers)
  - [Contribute](#contribute)
  - [License](#license)

## Installation

### Using Homebrew

```bash
brew tap kishaningithub/tap
brew install dynamodb-backup-restore
```

### Using Binary

```bash
# All unix environments with curl
curl -sfL https://raw.githubusercontent.com/kishaningithub/dynamodb-backup-restore/master/install.sh | sh -s -- -b /usr/local/bin

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/kishaningithub/dynamodb-backup-restore/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

## Example

### Backup single table

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -t employee-details -m backup -o backup-file
```

### Backup tables using regex pattern

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -p '.*-details' -m backup -o backup-file
```

### Backup all tables

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -p '.*' -m backup -o backup-file
```

### Restore

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -m restore -i backup-file
```

## Usage

```bash
dynamodb-backup-restore -h
Usage:
  dynamodb-backup-restore [OPTIONS]

Application Options:
  -t, --table-name=         Table name
  -p, --table-name-pattern= Table name pattern
  -m, --mode=               Mode of operation (backup,restore)
  -i, --input-backup-file=  Input backup file path
  -o, --output-backup-file= Output backup file path
  -e, --endpoint-url=       Endpoint url of destination dynamodb instance (Very useful for operating with local dynamodb instance)

Help Options:
  -h, --help                Show this help message
```

## Maintainers

- [@kishaningithub](https://github.com/kishaningithub)

## Contribute

1. Fork and fix/implement in a branch.
2. Make sure tests pass.
3. Make sure you've added new coverage.
4. Submit a PR.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2019 Kishan B
