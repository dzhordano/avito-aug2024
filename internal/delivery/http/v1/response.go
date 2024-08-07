package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

type authTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

type UserIdResponse struct {
	UserId string `json:"user_id"`
}

type response struct {
	Message string `json:"message"`
}

// messageResponse is a custom response type.
//
// swagger:response messageResponse
func messageResponse(c *gin.Context, statusCode int, message string) {
	// in: body
	fmt.Println(fmt.Errorf(message))
	c.AbortWithStatusJSON(statusCode, response{message})
}
