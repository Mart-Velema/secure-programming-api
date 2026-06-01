const SteamUser = require("steam-user");

const client = new SteamUser();

function getSteamClientStatus() {
  return {
    clientCreated: !!client,
    loggedOn: client.steamID !== null,

    credentialsConfigured:
      !!process.env.STEAM_USERNAME &&
      !!process.env.STEAM_PASSWORD,

    sharedSecretConfigured:
      !!process.env.STEAM_SHARED_SECRET,

    identitySecretConfigured:
      !!process.env.STEAM_IDENTITY_SECRET,
  };
}

module.exports = {
  client,
  getSteamClientStatus,
};