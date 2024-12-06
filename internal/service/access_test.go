package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/internal/service"
	"github.com/stretchr/testify/assert"
)

var (
	accessSecret = os.Getenv("ACCESS_SECRET")
)

func TestAccessService_Create(t *testing.T) {
	s := service.NewAccessService()

	type args struct {
		ctx context.Context
		td  *service.TokensData
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		validateFn func(t *testing.T, token string)
	}{
		{
			name: "successfully create token",
			args: args{
				ctx: context.Background(),
				td: &service.TokensData{
					IP:      "192.168.1.1",
					UserID:  "12345",
					TokenID: "token123",
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, token string) {
				parsedToken, err := jwt.ParseWithClaims(token, &service.TokenAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(accessSecret), nil
				})
				assert.NoError(t, err, "token parsing should not fail")

				if claims, ok := parsedToken.Claims.(*service.TokenAccessClaims); ok && parsedToken.Valid {
					assert.Equal(t, "192.168.1.1", claims.IP)
					assert.Equal(t, "12345", claims.Subject)
					assert.Equal(t, "tokens_service", claims.Issuer)
					assert.Equal(t, "token123", claims.Id)
				} else {
					t.Error("invalid token claims")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(tt.args.ctx, tt.args.td)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validateFn(t, got)
		})
	}
}
