package models

type AuthBasic struct {
	UserID string `json:"user_id"`
	APIKey string `json:"api_key"`
}
