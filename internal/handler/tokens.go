package handler

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/jwtGO/internal/service"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) createTokens(c *gin.Context) {
	userID := c.Query("user_id")
	if !isValidUUID(c, userID) {
		return
	}

	tokensData := &service.TokensData{
		UserID: userID,
		TokenID: generateTokenID(),
		IP: getIP(c),
	}

	accessToken, err := h.serv.CreateAccessToken(c, tokensData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating access token")
		return
	}

	refreshToken, err := h.serv.CreateRefreshToken(c, tokensData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func (h *Handler) refreshTokens(c *gin.Context) {
	userID := c.Query("user_id")
	if !isValidUUID(c, userID) {
		return
	}

	refreshTokenQuery := c.Query("refresh_token")
	if refreshTokenQuery == "" {
		newErrorResponse(c, http.StatusBadRequest, "Missing refresh_token")
		return
	}

	wasID, wasIP, err := h.serv.ParseRefreshToken(c, userID, refreshTokenQuery)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	nowIP := getIP(c)
	if nowIP != wasIP {
		// логика отправки сообщения пользователю о смене ip
	}

	tokensData := &service.TokensData{
		UserID: userID,
		TokenID: generateTokenID(),
		IP: nowIP,
	}

	accessToken, err := h.serv.CreateAccessToken(c, tokensData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating access token")
		return
	}

	refreshToken, err := h.serv.UpdateRefreshToken(c, wasID, tokensData)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func isValidUUID(c *gin.Context, userID string) bool {
    re := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	if !re.MatchString(userID) {
		if userID == "" {
			newErrorResponse(c, http.StatusBadRequest, "missing user_id")
			return false
		}

		newErrorResponse(c, http.StatusBadRequest, "invalid user_id")
		return false
	}

    return true
}

func getIP(c *gin.Context) string {
	r := c.Request
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	ip = r.RemoteAddr
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip
	}
	return host
}

func generateTokenID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
