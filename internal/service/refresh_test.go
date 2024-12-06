package service

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRefreshToken(t *testing.T) {
	type args struct {
		td *TokensData
	}
	tests := []struct {
		name string
		args
		wantErr bool
		validateFn func(t *testing.T, token string)
	}{
		{
			name: "Succesfully generate token",
			args: args{
				td: &TokensData{
					UserID: "test_user_id",
					TokenID: "test_token_id",
					IP: "111.111.1.1",
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &TokenRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(refreshSecret), nil
				})
				assert.NoError(t, err, "token parsing should not fail")

				if claims, ok := parsedToken.Claims.(*TokenRefreshClaims); ok && parsedToken.Valid {
					assert.Equal(t, "test_user_id", claims.Subject)
					assert.Equal(t, "tokens_service", claims.Issuer)
					assert.Equal(t, "test_token_id", claims.Id)
				} else {
					t.Error("invalid token claims")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRefreshToken(tt.args.td)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validateFn(t, got)
		})
	}
}
