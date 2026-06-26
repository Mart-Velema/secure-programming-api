# GuineaTrade API

<img alt="GuineaTrade logo" src="data:image/png;base64, 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">

The API for GuineaTrade

## Server Development

<hr>

### Installing prerequisites

To start working on the API, you must first install the Go language from [go.dev](https://go.dev/dl/) or run the commands below.

Installing `Go 1.26.4`.
```shell
cd
wget https://go.dev/dl/go1.26.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.4.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source $HOME/.profile
go version
```


The minimum required version is Go version `1.26.4`

You might also need to install [SQLite3](https://sqlite.org/download.html) as well. 
This can be either done through downloading and installing an installer from the website, or by using your favourite package manager of choice:

```shell
# Windows (Outdated, requires gcc)
winget install -e --id SQLite.SQLite

# Debian
sudo apt install sqlite3 gcc

# Fedora
sudo dnf install sqlite3 gcc
```

Once you've installed Go, you can run the following command to get started with the dependencies:
```shell
go mod tidy
```

_note_: Make sure you've also configured the `.env` file before you run the API, else it will *not* start

To run the API, use:
```shell
CGO_ENABLED=1  ENV_FILE="/path/to/env/.env" go run guineatrade.nhlstenden.com/src
```

<hr>

### .env file

The environment variables are stored in a `.env` file in the project root. It contains the following entries:

```shell
PORT="3000" # The port to run the API on

# Database
ENCRYPTION_KEY="Fy2mcmCCOS6LsLRpkoeJSFCNVaixHN9oLEdoU+6lu2A=" # A Base64 encoded 32 byte sequence used as the master key to encrypting and decrypting the database entries
SQLITE_FILE_LOCATION="./guineatrade.db" # The filepath where the SQLite3 database file is located

# JWT
JWT_SECRET_KEY="superlongkey" # A long string of random data to act as the signing key for JWT. Recommended to have at least 64 characters or more
JWT_TIMEOUT_MINUTES=3 # Time a JWT is valid in minutes
JWT_REFRESH_DAYS=7 # Time the refresh token is valid in days

# Backpack.tf
BACKPACK_API_KEY="12345678901234567890abcd" # The Backpack.TF API key
BACKPACK_API_HASH="sha256:39a999a6d0aad5c4be9ea3c952dd6331d6c14ff2b2c0f1e1e99fb11e8653e78f" # The backpack.tf API hash
ITEM_CONSTANTS="/path/to/item-constants.json" # The absolute filepath to the item-constants.json

# Steam API
STEAM_API_HASH="sha256:STEAM_API_HASH_HERE" # The Steam API hash

# Steam Bot
BOT_PORT="3001" # Internal Steam bot port
STEAM_BOT_URL=http://steam-bot:3001 # Internal URL used by the Go API to communicate with the Steam bot
BOT_API_KEY=super-secret-key # Shared API key used between the Go API and Steam bot

STEAM_USERNAME=guineatradebot # Steam account username
STEAM_PASSWORD= # Steam account password
STEAM_SHARED_SECRET= # Optional. Enables automatic Steam Guard login
STEAM_IDENTITY_SECRET= # Optional. Enables automatic trade confirmations

```

API hashes can be generated with the following OpenSSL pipeline:
```sh
openssl s_client -connect backpack.tf:443 -servername backpack.tf | openssl x509 -pubkey -noout | openssl pkey -pubin -outform der | openssl dgst -sha256
```

In case that the binary or project is executed in a different directory than the .env file, you can supply the filepath using the `--env` CLI flag:

```
./gt --env /path/to/.env
```

<hr>

## Steam Bot

The Steam trading bot runs as a separate internal service and is started automatically through Docker Compose together with the API.

### Logging Into Steam

If `STEAM_SHARED_SECRET` is not configured, the bot requires a Steam Guard code.

Login through the API:

```http
POST /api/v1/steam/login
Authorization: Basic <credentials>
Content-Type: application/json
```

Request body:

```json
{
  "authCode": "ABCDE"
}
```

Replace `ABCDE` with the current Steam Guard code from the Steam Mobile App.

Verify login:

```http
GET /api/v1/steam/status
Authorization: Bearer <JWT>
```

Successful login returns:

```json
{
  "loggedOn": true,
  "steamId": "7656119..."
}
```


### Database

The database is a generic SQLite3 database. it is managed by [GORM](https://gorm.io/docs/), which automatically adds in the required columns and constraints, sets-up encryption for the fields that require it and manage the relations between the tables.
To assess the database instance in the code, you can use the database singleton:

```go
db := database.getInstance()
```

The database can also be seeded with random values. You can enable seeding by passing `--seed` as a CLI argument

_note_: Seeding the database permanently and destructively deletes the entire database. Really, it literally deletes the file permanently.

### Middlewares

There are two different middleware suites available: JWT and TOTP.

#### JWT

JSON Web Token is used for basic authentication and verifying that a user is not a web crawler. JWT's are send and checked using the Authorization HTTP header. Each user gets a unique JWT, and are valid for only a few minutes.
With `middleware.ExtractTokenUser()`, you can get the User from the current context.

```http request
GET /api/v1/auth/me HTTP/1.1
Authorization: Bearer <JWT>
```

#### TOTP

Time based One Time Passwords are 6 digit codes that are valid for 30 seconds. TOTP should be used for sensitive transactions which could cost us or the user money.
The user may send a recovery code in certain contexts to reset the TOTP. 

_the recovery code should never be used for verification_

```http request
POST /api/v1/auth/mfa/totp/reset HTTP/1.1
Authorization: Bearer <JWT>
X-TOTP-Token: <TOTP>
X-Recovery-Code: <Recovery>
```

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