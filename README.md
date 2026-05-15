# GuineaTrade API

<img alt="GuineaTrade logo" src="data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAKAAAACgCAYAAACLz2ctAAAACXBIWXMAAC4jAAAuIwF4pT92AAAF4klEQVR4Xu2du4pUQRRFe0QDHwMaiIYaCJMJgiAiRoKR+Bdm/oGhfoHZ/IUYCUYiIgiCmWBiqJGCr8BAbTEw63UuZ0/fxzJ1975VqxYl3LK6d1bL+vNrAtPdmcAY24Z4qK3JIgkMIKCAA6D5kT4CCtjH0qYBBBRwADQ/0kdAAftY2jSAgAIOgOZH+ggoYB9LmwYQUMAB0PxIHwEF7GNp0wACU3jr3np68evJAEoH+JGdm+0PG/UauwO2r7eFFQIKWKFltp2AArYjtbBCQAErtMy2E1DAdqQWVggoYIWW2XYCCtiO1MIKAQWs0DLbTkAB25FaWCGQeEu+lZOLp1/uVea9Mfv568bI38DJEyzXnbqxex9Vjv1kxR0QLaOhFAEFTJG1FxFQQITJUIqAAqbI2osIKCDCZChFQAFTZO1FBBQQYTKUIqCAKbL2IgIKiDAZShE4nCre1LutuxljP+HYxK3695Rz4MQEDdUdEGEylCKggCmy9iICCogwGUoRUMAUWXsRAQVEmAylCChgiqy9iIACIkyGUgQUMEXWXkRAAREmQykClTsh6K4HffPefYcjBWgpvYE7Jsgtd8ClGDbSeSrgSBdmKcNSwKWs9EjnqYAjXZilDEsBl7LSI52nAo50YZYyLAVcykqPdJ4KONKFWcqwFHApKz3SebbfCdl/fQtO9Q3KXb/wFuWevdtDOdqHygKh7nnQvv0VXbfHrbN2B2zFaVmVgAJWiZlvJaCArTgtqxJQwCox860EFLAVp2VVAgpYJWa+lYACtuK0rEpAAavEzLcSUMBWnJZVCaxPQlrveuy/ZkOgJxL0Tf6dq71v6D9+O8smAlNnjn9Ayb3T71Bu/wU7uejmTO/8/Pm2LeLVjjsgWm5DKQIKmCJrLyKggAiToRQBBUyRtRcRUECEyVCKgAKmyNqLCCggwmQoRUABU2TtRQQUEGEylCLQfickNdCu3pc/H6Kq86sHKEdD9LlXjtyllbPIuQPOYhmnOwkFnO7azWLkCjiLZZzuJBRwums3i5Er4CyWcbqTUMDprt0sRq6As1jG6U5CAae7drMYuQLOYhmnO4n1bzmQ/7u/oncBKAr6LVr0TsOpY1/oo2eR+/R9F80D36m51HunBv4Cu3dC0CoaihHwn+AYWosJAQUklMzECChgDK3FhIACEkpmYgQUMIbWYkJAAQklMzECChhDazEhoICEkpkYgdHfCVnaCQdd6blwcQekK24uQkABI1gtpQQUkJIyFyGggBGsllICCkhJmYsQUMAIVkspAQWkpMxFCChgBKullIACUlLmIgTWJyHreyEb/8DffWi/O7JxYMXAmXPXip842PjH988P9oHFp8G7HutW5JU7YHEBjPcSUMBenrYVCShgEZjxXgIK2MvTtiIBBSwCM95LQAF7edpWJKCARWDGewkoYC9P24oEFLAIzHgvAfS2+t8jR/0tWhTL7YuvULT7xISecDx6cxmNj4bubOdbr9bDQ265A9KVNBchoIARrJZSAgpISZmLEFDACFZLKQEFpKTMRQgoYASrpZSAAlJS5iIEFDCC1VJKQAEpKXMRAuhtdfHJWzkxoWOkv09C+7pz3ScXdHzddz3oc90BKSlzEQIKGMFqKSWggJSUuQgBBYxgtZQSUEBKylyEgAJGsFpKCSggJWUuQkABI1gtpQQUkJIyFyGwtZMQOpvuX2qnz23PHYWNP2AOxgonHLCR3fWgZe6AlJS5CAEFjGC1lBJQQErKXISAAkawWkoJKCAlZS5CQAEjWC2lBBSQkjIXIaCAEayWUgIKSEmZixBInIR0DxTdMaEPHfvJythPLihnmnMHpKTMRQgoYASrpZSAAlJS5iIEFDCC1VJKQAEpKXMRAgoYwWopJaCAlJS5CAEFjGC1lBJQQErKXITAFE5COifeeqrSObD/uha1Ju6AIYusZQQUkHEyFSKggCGw1jICCsg4mQoRUMAQWGsZAQVknEyFCChgCKy1jIACMk6mQgQUMATWWkbgN6sBbrPn25MFAAAAAElFTkSuQmCC">

The API for GuineaTrade

## Server Development

<hr>

### Installing prerequisites

To start working on the API, you must first install the Go language from [go.dev](https://go.dev/dl/).
The minimum required version is Go version `1.26.2`

You might also need to install [SQLite3](https://sqlite.org/download.html) as well. 
This can be either done through downloading and installing an installer from the website, or by using your favourite package manager of choice:

```shell
# Windows
winget install -e --id SQLite.SQLite

# Debian
sudo apt install sqlite3

# Fedora
sudo dnf install sqlite3
```

Once you've installed Go, you can run the following command to get started with the dependencies:
```shell
go mod tidy
```

_note_: Make sure you've also configured the `.env` file before you run the API, else it will *not* start

To run the API, use:
```shell
go run guineatrade.nhlstenden.com/src
```

<hr>

### .env file

The environment variables are stored in a `.env` file in the project root. It contains the following entries:

```shell
PORT="3000" # The port to run the API on

# Database
ENCRYPTION_KEY="Fy2mcmCCOS6LsLRpkoeJSFCNVaixHN9oLEdoU+6lu2A=" # A Base64 encoded 32 byte sequence used as the master key to encrypting and decrypting the database entries
ENCRYPTION_PASSCODE="My Super Secure Password!" # A Generic any-length string that will be used as the master password. Overrides the value of ENCRYPTION_KEY
SQLITE_FILE_LOCATION="./guineatrade.db" # The filepath where the SQLite3 database file is located

```

<hr>

### Database

The database is a generic SQLite3 database. it is managed by [GORM](https://gorm.io/docs/), which automatically adds in the required columns and constraints, sets-up encryption for the fields that require it and manage the relations between the tables.
To assess the database instance in the code, you can use the database singleton:

```go
db := database.getInstance()
```

The database can also be seeded with random values. You can enable seeding by passing `--seed` as a CLI argument

_note_: Seeding the database permanently and destructively deletes the entire database. Really, it literally deletes the file permanently.

## Server Hosting

The API has a docker setup included. It can be run using the docker-compose tool:

```shell
docker-compose build
docker-compose up -d
```
_note_: It is recommended to set `SQLITE_FILE_LOCATION` to a file placed within the `/data` directory when using Docker, for example:
```shell
SQLITE_FILE_LOCATION="/data/gt.db"
```