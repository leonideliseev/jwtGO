package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) createTokens(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		newErrorResponse(c, http.StatusBadRequest, "missing user id")
		return
	}
	ip := getIP(c)

	accessToken, err := h.serv.GenerateAccessToken(c, userID, ip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating access token")
		return
	}

	refreshToken, err := h.serv.GenerateRefreshToken(c, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	response := TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func (h *Handler) refreshTokens(c *gin.Context) {
	userID := c.Query("user_id")
	refreshTokenQuery := c.Query("refresh_token")
	if userID == "" || refreshTokenQuery == "" {
		newErrorResponse(c, http.StatusBadRequest, "Missing user_id or refresh_token")
		return
	}

	err := h.serv.CheckRefreshToken(c, userID, refreshTokenQuery)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ip := getIP(c)

	accessToken, err := h.serv.GenerateAccessToken(c, userID, ip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating access token")
		return
	}

	refreshToken, err := h.serv.UpdateRefreshToken(c, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	response := TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
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
