# Bridge Tab

## Overview

Bridge Tab is a tool designed to manage duplicate bridge tournaments. It provides functionalities for organizers or umpires to prepare and manage tournaments, check scores, and manage users.
It also provides a http server that allows contestants to record scores for rounds they play in, a http server that allows administrating a tournament and a http server for tablet kiosk setup per table.


**Notice** This project is still work in progress.

## Features

- Manage duplicate bridge tournaments
- Play registration by contestants
- Play summary in CSV
- Calculating scores in Matchpoints
- Tournament schedule
- Frontend for tables
- Admin panel
- CLI

## Current roadmap

- [ ] Collecting cards per hand for a board
- [ ] Creating prettier frontends
- [ ] Implementing IMPs
- [ ] More QoL features
- [ ] Boards analysis engine

## Getting Started

### Prerequisites

- Go (latest version)
- Postgres

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

Optional `HTTP_BASE_PATH` (for example `/app`) prefixes redirects, links, and cookies when this frontend is served behind a path-based reverse proxy. Leave unset when the app is at the domain root. Nginx should strip the prefix before proxying (e.g. `location /app/ { proxy_pass http://127.0.0.1:3000/; }`).

#### Admin

:warning: This is a very early stage. It was done with a lot of AI and is more of a concept/prototype.

The admin panel is a separate server. It needs the same database and a long random password:

```bash
export DATABASE_STRING=<your string here>
export ADMIN_PASSWORD=<long random password>
```

`ADMIN_PORT` defaults to `3001`. Set `ADMIN_COOKIE_SECURE=1` when the panel is served over HTTPS.

Optional `ADMIN_BASE_PATH` (for example `/admin`) prefixes redirects, links, and cookies for path-based reverse proxies. Nginx should strip the prefix before proxying.

```bash
make admin
```

To only build it:
```bash
make build-admin
```

#### Tablet

The table-tablet frontend is a separate server for devices that stay at a table and record results.

```bash
export DATABASE_STRING=<your string here>
export TABLET_PASSWORD=<long random password>
```

`TABLET_PORT` defaults to `3002`. Set `TABLET_COOKIE_SECURE=1` when the panel is served over HTTPS.

Optional `TABLET_BASE_PATH` (for example `/tablet`) prefixes redirects, links, and cookies for path-based reverse proxies. Nginx should strip the prefix before proxying.

```bash
make tablet
```

To only build it:

```bash
make build-tablet
```

After login, each device is assigned a tournament and table. Once that tournament is live, the assignment cannot be changed on that device. A replacement device can still be set up.

#### CLI

Run the CLI tool with the following command:
```bash
make build-cli
```

Set up the database connection string:
```bash
EXPORT DATABASE_STRING=<your string here>
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
