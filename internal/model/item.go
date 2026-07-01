package model

import "time"

type ItemStatus string

const (
	StatusActive   ItemStatus = "active"
	StatusOffline  ItemStatus = "offline"
	StatusTraded   ItemStatus = "traded"
	StatusPending  ItemStatus = "pending"
)

type Item struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	City            string     `json:"city"`
	Condition       string     `json:"condition"`
	Publisher       string     `json:"publisher"`
	DesiredCategory string     `json:"desiredCategory"`
	Description     string     `json:"description"`
	ViewCount       int        `json:"viewCount"`
	FavoriteCount   int        `json:"favoriteCount"`
	TradeIntentCount int       `json:"tradeIntentCount"`
	Status          ItemStatus `json:"status"`
	IsFavorited     bool       `json:"isFavorited"`
	HasCommunicated bool       `json:"hasCommunicated"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type ItemCreateRequest struct {
	Title           string `json:"title"`
	Category        string `json:"category"`
	City            string `json:"city"`
	Condition       string `json:"condition"`
	Publisher       string `json:"publisher"`
	DesiredCategory string `json:"desiredCategory"`
	Description     string `json:"description"`
}

type ItemUpdateRequest struct {
	Title           string `json:"title"`
	Category        string `json:"category"`
	City            string `json:"city"`
	Condition       string `json:"condition"`
	DesiredCategory string `json:"desiredCategory"`
	Description     string `json:"description"`
}

type ItemQuery struct {
	Keyword  string     `form:"keyword"`
	Category string     `form:"category"`
	City     string     `form:"city"`
	Status   ItemStatus `form:"status"`
}

type Statistics struct {
	TotalItems     int `json:"totalItems"`
	ActiveItems    int `json:"activeItems"`
	OfflineItems   int `json:"offlineItems"`
	TradedItems    int `json:"tradedItems"`
	TotalViews     int `json:"totalViews"`
	TotalFavorites int `json:"totalFavorites"`
	TotalIntents   int `json:"totalIntents"`
}
