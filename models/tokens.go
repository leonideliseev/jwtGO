package models

type Refresh struct {
	TokenID          string `db:"token_id"`
	IP               string `db:"ip"`
	RefreshTokenHash string `db:"refresh_token_hash"`
}
