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
  - [Install](#install)
    - [Using Homebrew](#using-homebrew)
    - [Using Binary](#using-binary)
  - [Example](#example)
  - [Usage](#usage)
  - [Maintainers](#maintainers)
  - [Contribute](#contribute)
  - [License](#license)

## Install

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

### Backup

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -t employee-details  -m backup -o employee-details.json
```

### Restore

```bash
AWS_REGION=eu-west-1 AWS_SDK_LOAD_CONFIG=true dynamodb-backup-restore -t employee-details  -m restore -i employee-details.json
```

### Backing up multiple tables at once

```bash
export AWS_REGION=eu-west-1 
export AWS_SDK_LOAD_CONFIG=true
cat tables | xargs -I {} dynamodb-backup-restore -t {} -m backup -o {}.json
```

### Restoring up multiple tables at once

```bash
export AWS_REGION=eu-west-1 
export AWS_SDK_LOAD_CONFIG=true
cat tables | xargs -I {} dynamodb-backup-restore -t {}  -m restore -i {}.json
```

## Usage

```bash
dynamodb-backup-restore -h
Usage:
  main [OPTIONS]

Application Options:
  -t, --table-name= Name of the dynamo db table
  -m, --mode=       Mode of operation (backup,restore)
  -o, --output=     Output file for backup
  -i, --input=      Input file for restore

Help Options:
  -h, --help        Show this help message
```

## Maintainers

- [@kishaningithub](https://github.com/kishaningithub)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2019 Kishan B
