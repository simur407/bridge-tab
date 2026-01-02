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

- [ ] More tests, especially integration/e2e
- [ ] Adding better frontend for contestants
- [ ] Tournament scoring

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

#### CLI

Run the CLI tool with the following command:
```bash
make build-cli
```

Then you can run it with the following command:
```bash
./build/bridge-tab --help
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the GNU Affero General Public License v3.0. See the LICENSE file for details.
