# GoBank

GoBank is a simple banking application built with Go.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Introduction

GoBank is a Go-based backend application that simulates basic banking operations. It includes features like creating accounts, transferring money, and listing accounts.

## Features

- Account creation with user details
- Money transfer between accounts
- Listing all accounts

## Getting Started

### Prerequisites

To run this project, you need to have Go installed on your machine. Make sure you also have a PostgreSQL database available.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mounis-bhat/go-bank.git
   cd gobank
   Copy the .env.example file to .env and configure the database connection details.
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Build the application:

   ```bash
    make build
   ```

4. Run the application:

   ```bash
   make run
   ```

5. Visit http://localhost:8080 to access the API.

6. To run the tests:

   ```bash
   make test
   ```

## Usage

Use the following endpoints to access the API:

| Endpoint  | Method | Description                     |
| --------- | ------ | ------------------------------- |
| /account  | GET    | Get account details             |
| /account  | PATCH  | Update account details          |
| /account  | POST   | Create a new account            |
| /account  | DELETE | Delete an account               |
| /accounts | GET    | List all accounts               |
| /transfer | POST   | Transfer money between accounts |
| /login    | POST   | Login to an account             |

## Testing

The tests are written using the standard Go testing package. To run the tests, use the following command:

```bash
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
