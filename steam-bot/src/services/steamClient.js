const SteamUser = require("steam-user");
const SteamCommunity = require("steamcommunity");
const TradeOfferManager = require("steam-tradeoffer-manager");

const client = new SteamUser();
const community = new SteamCommunity();

const manager = new TradeOfferManager({
  steam: client,
  community,
  language: "en",
});

let lastError = null;
let isLoggingIn = false;

client.on("loggedOn", () => {
  console.log("Steam bot logged in successfully");
  lastError = null;
  isLoggingIn = false;
  client.setPersona(SteamUser.EPersonaState.Online);
});

client.on("webSession", (sessionID, cookies) => {
  console.log("Steam web session established");

  community.setCookies(cookies);

  manager.setCookies(cookies, (err) => {
    if (err) {
      console.error("Failed to set manager cookies:", err);
      return;
    }

    console.log("Trade manager ready");
  });
});

client.on("error", (err) => {
  console.error("Steam client error:", err.message);
  lastError = err.message;
  isLoggingIn = false;
});

client.on("disconnected", (eresult, msg) => {
  console.log("Steam bot disconnected:", eresult, msg);
});

function getSteamClientStatus() {
  return {
    clientCreated: !!client,
    loggedOn: client.steamID !== null,
    isLoggingIn,
    credentialsConfigured:
      !!process.env.STEAM_USERNAME && !!process.env.STEAM_PASSWORD,
    sharedSecretConfigured: !!process.env.STEAM_SHARED_SECRET,
    identitySecretConfigured: !!process.env.STEAM_IDENTITY_SECRET,
    steamId: client.steamID ? client.steamID.getSteamID64() : null,
    lastError,
  };
}

function loginToSteam(authCode) {
  if (client.steamID) {
    return {
      ok: true,
      message: "Already logged in",
    };
  }

  if (isLoggingIn) {
    return {
      ok: false,
      message: "Login already in progress",
    };
  }

  if (!process.env.STEAM_USERNAME || !process.env.STEAM_PASSWORD) {
    return {
      ok: false,
      message: "Missing STEAM_USERNAME or STEAM_PASSWORD",
    };
  }

  isLoggingIn = true;
  lastError = null;

  const logOnOptions = {
    accountName: process.env.STEAM_USERNAME,
    password: process.env.STEAM_PASSWORD,
  };

  if (authCode) {
    logOnOptions.authCode = authCode;
  }

  client.logOn(logOnOptions);

  return {
    ok: true,
    message: "Login attempt started",
  };
}

module.exports = {
  client,
  community,
  manager,
  getSteamClientStatus,
  loginToSteam,
};