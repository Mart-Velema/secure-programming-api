# Steam Bot

Steam trading bot service for GuineaTrade.

# IMPORTANT!!!
The steam bot is its own account, make sure you have the correct credentials before continuing

## Features

Currently implemented:

* Steam client initialization
* Steam login via username/password
* Steam Guard authentication
* Steam web session creation
* TradeOfferManager integration
* Inventory retrieval
* Status endpoints

## Installation

Install dependencies:

```bash
npm install
```

## Configuration

Create a `.env` file:

```env
BOT_PORT=3001

STEAM_USERNAME=
STEAM_PASSWORD=
STEAM_SHARED_SECRET=
STEAM_IDENTITY_SECRET=
```

### Environment Variables

| Variable                | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| `BOT_PORT`              | Port the bot service runs on                                 |
| `STEAM_USERNAME`        | Steam bot account username                                   |
| `STEAM_PASSWORD`        | Steam bot account password                                   |
| `STEAM_SHARED_SECRET`   | Used for automatic Steam Guard codes (not yet implemented)   |
| `STEAM_IDENTITY_SECRET` | Used for automatic trade confirmations (not yet implemented) |

## Starting the Bot

Run:

```bash
npm run dev
```

Expected output:

```txt
Steam bot service running on port 3001
```

## Login

Send a login request:

### Example: PowerShell

```powershell
Invoke-RestMethod -Uri "http://localhost:3001/steam/login" `
-Method Post `
-ContentType "application/json" `
-Body "{}"
```

### Steam Guard

Steam Guard is enabled, enter the code from the Steam Mobile App when prompted.

Successful login:

```txt
Steam bot logged in successfully
Steam web session established
Trade manager ready
```

## Verify Status

Open:

```txt
http://localhost:3001/steam/status
```

Expected response:

```json
{
  "loggedOn": true
}
```

## Inventory

Open:

```txt
http://localhost:3001/steam/inventory
```

Default values:

| Game             | App ID | Context ID |
| ---------------- | ------ | ---------- |
| Team Fortress 2  | 440    | 2          |
| Counter-Strike 2 | 730    | 2          |

## Available Endpoints

### Health Check

```http
GET /health
```

### Configuration Status

```http
GET /config
```

### Steam Status

```http
GET /steam/status
```

### Steam Login

```http
POST /steam/login
```

### Bot Inventory

```http
GET /steam/inventory
```

## Current Status

Implemented:

* Steam authentication
* Steam Guard login
* Inventory retrieval
* TradeOfferManager setup

Not yet implemented:

* Automatic Steam Guard login (`shared_secret`)
* Automatic trade confirmations (`identity_secret`)
* Send trade offers
* Accept trade offers
* Cancel trade offers
* Trade status monitoring
* Integration with Go API