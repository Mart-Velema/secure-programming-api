# Steam Bot

Steam trading bot service for GuineaTrade.

# IMPORTANT

The Steam bot uses its own Steam account. Make sure you are using the correct Steam account credentials before continuing.

## Features

Currently implemented:

* Steam client initialization
* Steam login via username/password
* Steam Guard authentication
* Steam web session creation
* TradeOfferManager integration
* Inventory retrieval
* Trade URL validation
* Create trade offers
* Accept incoming trade offers
* Cancel outgoing trade offers
* View trade offer details
* List active trade offers
* View trade offer history
* Internal API key protection

---

## Installation

Install dependencies:

```bash
npm install
```

---

## Configuration

Create a `.env` file:

```env
BOT_PORT=3001

STEAM_USERNAME=
STEAM_PASSWORD=
STEAM_SHARED_SECRET=
STEAM_IDENTITY_SECRET=

BOT_API_KEY=
```

### Environment Variables

| Variable                | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| `BOT_PORT`              | Port the bot service runs on                                 |
| `STEAM_USERNAME`        | Steam bot account username                                   |
| `STEAM_PASSWORD`        | Steam bot account password                                   |
| `STEAM_SHARED_SECRET`   | Used for automatic Steam Guard codes (not yet implemented)   |
| `STEAM_IDENTITY_SECRET` | Used for automatic trade confirmations (not yet implemented) |
| `BOT_API_KEY`           | Internal API key required for protected endpoints            |

---

## Starting the Bot

Run:

```bash
npm run dev
```

Expected output:

```txt
Steam bot service running on port 3001
```

---

## Authentication

All endpoints except `/health` require the API key.

Example header:

```http
X-API-Key: your-secret-api-key
```

### PowerShell Example

```powershell
Invoke-RestMethod `
  -Uri "http://localhost:3001/steam/status" `
  -Headers @{ "X-API-Key" = "your-secret-api-key" }
```

---

## Login

Send a login request:

```powershell
Invoke-RestMethod `
  -Uri "http://localhost:3001/steam/login" `
  -Method Post `
  -ContentType "application/json" `
  -Headers @{ "X-API-Key" = "your-secret-api-key" } `
  -Body "{}"
```

### Steam Guard

Steam Guard is enabled on the bot account.

When prompted in the console, enter the Steam Guard code from the Steam Mobile App.

Successful login:

```txt
Steam bot logged in successfully
Steam web session established
Trade manager ready
```

---

## Verify Status

```powershell
Invoke-RestMethod `
  -Uri "http://localhost:3001/steam/status" `
  -Headers @{ "X-API-Key" = "your-secret-api-key" }
```

Expected response:

```json
{
  "loggedOn": true
}
```

---

## Inventory

Example:

```txt
GET /steam/inventory?appId=440&contextId=2
```

## Available Endpoints

### Public

| Method | Endpoint  | Description  |
| ------ | --------- | ------------ |
| GET    | `/health` | Health check |

### Protected

| Method | Endpoint                                   | Description                            |
| ------ | ------------------------------------------ | -------------------------------------- |
| GET    | `/config`                                  | Configuration status                   |
| GET    | `/steam/status`                            | Steam client status                    |
| POST   | `/steam/login`                             | Login to Steam                         |
| GET    | `/steam/inventory`                         | Retrieve inventory                     |
| POST   | `/steam/trade-url/validate`                | Validate trade URL                     |
| POST   | `/steam/trade-offers/dry-run`              | Validate trade request without sending |
| POST   | `/steam/trade-offers`                      | Create and send trade offer            |
| GET    | `/steam/trade-offers`                      | List active trade offers               |
| GET    | `/steam/trade-offers/history`              | List historical trade offers           |
| GET    | `/steam/trade-offers/:tradeOfferId`        | Retrieve trade offer details           |
| POST   | `/steam/trade-offers/:tradeOfferId/cancel` | Cancel outgoing trade offer            |
| POST   | `/steam/trade-offers/:tradeOfferId/accept` | Accept incoming trade offer            |

---

## Trade Offer Request Example

```json
{
  "tradeUrl": "https://steamcommunity.com/tradeoffer/new/?partner=123456789&token=abcdef",
  "itemsToGive": [
    {
      "appId": 440,
      "contextId": 2,
      "assetId": "20540621909"
    }
  ],
  "message": "GuineaTrade trade offer"
}
```

---

## Current Status

Implemented:

* Steam authentication
* Steam Guard login
* Steam inventory retrieval
* TradeOfferManager setup
* Trade URL validation
* Send trade offers
* Accept incoming trade offers
* Cancel trade offers
* View trade offer details
* Active trade offer listing
* Trade offer history
* Internal API authentication
* Automatic Steam Guard login (`STEAM_SHARED_SECRET`) (optional)
* Automatic trade confirmations (`STEAM_IDENTITY_SECRET`) (optional)

Not yet implemented:

* Trade status polling
* Event notifications/webhooks
* Integration with Go API
* Persistent trade storage
* Order synchronization