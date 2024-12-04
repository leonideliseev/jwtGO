package models

type Refresh struct {
	UserID string `db:"user_id"`
	IP string `db:"ip"`
	RefreshTokenHash string `db:"refresh_token_hash"`
} 
