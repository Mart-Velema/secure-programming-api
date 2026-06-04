package steam

import "time"

type SteamInventoryResponse struct {
	AppID     int         `json:"appId"`
	ContextID string      `json:"contextId"`
	Count     int         `json:"count"`
	Inventory []SteamItem `json:"inventory"`
}

type SteamItem struct {
	AssetID        string `json:"assetId"`
	MarketHashName string `json:"marketHashName"`
	Name           string `json:"name"`
	Tradable       bool   `json:"tradable"`
	Marketable     bool   `json:"marketable"`
}

type SendTradeOfferRequest struct {
	TradeURL    string             `json:"tradeUrl"`
	ItemsToGive []TradeOfferItem   `json:"itemsToGive"`
	Message     string             `json:"message"`
}

type TradeOfferItem struct {
	AppID     int    `json:"appId"`
	ContextID string    `json:"contextId"`
	AssetID   string `json:"assetId"`
}

type SendTradeOfferResponse struct {
	OK           bool   `json:"ok"`
	TradeOfferID string `json:"tradeOfferId"`
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
}

type TradeOfferResponse struct {
	OK    bool       `json:"ok"`
	Offer TradeOffer `json:"offer"`
	Error string     `json:"error,omitempty"`
}

type TradeOffer struct {
	ID             string           `json:"id"`
	State          int              `json:"state"`
	StateName      string           `json:"stateName"`
	Partner        string           `json:"partner"`
	Message        string           `json:"message"`
	Created        time.Time        `json:"created"`
	Updated        time.Time        `json:"updated"`
	Expires        time.Time        `json:"expires"`
	ItemsToGive    []TradeOfferAsset `json:"itemsToGive"`
	ItemsToReceive []TradeOfferAsset `json:"itemsToReceive"`
}

type TradeOfferAsset struct {
	AssetID        string `json:"assetId"`
	AppID          int    `json:"appId"`
	ContextID      string `json:"contextId"`
	MarketHashName string `json:"marketHashName"`
}

type TradeOfferListResponse struct {
	OK       bool             `json:"ok"`
	Sent     []TradeOfferBrief `json:"sent"`
	Received []TradeOfferBrief `json:"received"`
	Error    string           `json:"error,omitempty"`
}

type TradeOfferBrief struct {
	ID        string    `json:"id"`
	State     int       `json:"state"`
	StateName string    `json:"stateName"`
	Partner   string    `json:"partner"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type BotStatusResponse struct {
	ClientCreated            bool    `json:"clientCreated"`
	LoggedOn                 bool    `json:"loggedOn"`
	IsLoggingIn              bool    `json:"isLoggingIn"`
	CredentialsConfigured    bool    `json:"credentialsConfigured"`
	SharedSecretConfigured   bool    `json:"sharedSecretConfigured"`
	IdentitySecretConfigured bool    `json:"identitySecretConfigured"`
	SteamID                  *string `json:"steamId"`
	LastError                *string `json:"lastError"`
}

type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type AcceptTradeOfferResponse struct {
	OK           bool   `json:"ok"`
	TradeOfferID string `json:"tradeOfferId"`
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
}

type CancelTradeOfferResponse struct {
	OK           bool   `json:"ok"`
	TradeOfferID string `json:"tradeOfferId"`
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
}