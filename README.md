# Bridge Tab

## Overview

Bridge Tab is a tool designed to manage duplicate bridge tournaments. It provides functionalities for organizers or umpires to prepare and manage tournaments, check scores, and manage users.
It also provides a http server that allows contestants to record scores for rounds they play in.

**Notice** This project is still work in progress.

## Features

- Manage duplicate bridge tournaments
- Round registration by contestants
- Rounds summary in CSV

## Current roadmap

- [x] More tests, especially integration/e2e
- [ ] Adding better frontend for contestants
- [x] Tournament scoring
- [ ] Admin panel

## Getting Started

### Prerequisites

- Go (latest version)

### Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:simur407/bridge-tab.git
   cd bridge-tab
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Usage

Set up the database connection string:
```bash
EXPORT DATABASE_STRING=<your string here>
```

#### HTTP

To build and run HTTP server use this command:
```bash
make http
```

To only build the HTTP server, go with the following command:
```bash
make build-http
```

Then you can run it with the following command:
```bash
make run-http
```

#### Admin

:warning: This is a very early stage. It was done with a lot of AI and is more of a concept/prototype.

The admin panel is a separate server. It needs the same database and a long random password:

```bash
export DATABASE_STRING=<your string here>
export ADMIN_PASSWORD=<long random password>
```

`ADMIN_PORT` defaults to `3001`. Set `ADMIN_COOKIE_SECURE=1` when the panel is served over HTTPS.

```bash
make admin
```

To only build it:
```bash
make build-admin
```

#### CLI

Run the CLI tool with the following command:
```bash
make build-cli
```

Then you can run it with the following command:
```bash
./build/bridge-tab --help
```

### Running tests

Run all Go tests with:
```bash
make test
```

Or equivalently:
```bash
cd backend && go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the GNU Affero General Public License v3.0. See the LICENSE file for details.

## AI

Some parts of this project were created with AI assistance. Generally AI assisted code is accepted, but it has to be thoroughly reviewed and up to codebase standards.
