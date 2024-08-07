package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"

	userTypeCtx = "user-type"
)

func (h *Handler) isAuthorized(c *gin.Context) {
	userType, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid auth token",
		})
	}

	c.Set(userTypeCtx, userType)
}

func (h *Handler) isModerator(c *gin.Context) {
	// Check if the auth_tokens is a moderator auth_tokens
	userType, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid auth token",
		})

		return
	}

	if userType != "moderator" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "only moderators are allowed",
		})
	}

	c.Set(userTypeCtx, userType)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	hParts := strings.Split(header, " ")
	if len(hParts) != 2 || hParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(hParts[1]) == 0 {
		return "", errors.New("auth_token is empty")
	}

	return h.tokensManager.Parse(hParts[1])
}
