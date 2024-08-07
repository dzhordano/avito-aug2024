package v1

import (
	"errors"

	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.GET("/dummyLogin", h.dummyLogin)

		auth.POST("/register", h.userRegister)
		auth.POST("/login", h.userLogin)
	}
}

// @Summary		Dummy Login
// @Description	get token corresponding to dummy user
// @ID				dummyLogin
// @Tags			auth
// @Produce		json
// @Param			userType	query		string	false	"userType"	Enums(client, moderator)
// @Success		200			{object}	authTokenResponse
// @Failure		400			{object}	response
// @Failure		500			{object}	response
// @Router			/auth/dummyLogin [get]
func (h *Handler) dummyLogin(c *gin.Context) {
	inp := c.Query("userType")
	if !domain.UserType(inp).Validate() {
		messageResponse(c, http.StatusBadRequest, "invalid user-type query")

		return
	}

	token, err := h.services.Users.DummyLogin(inp)
	if err != nil {
		messageResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, authTokenResponse{AuthToken: token})
}

// @Summary		User Register
// @Description	register new user with email, password and userType
// @ID				userRegister
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		dtos.UserRegisterInput	true	"User info"
// @Success		201		{object}	UserIdResponse
// @Failure		400		{object}	response
// @Failure		409		{object}	response
// @Failure		500		{object}	response
// @Router			/auth/register [post]
func (h *Handler) userRegister(c *gin.Context) {
	var inp dtos.UserRegisterInput
	if err := c.BindJSON(&inp); err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := inp.Validate(); err != nil {
		messageResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := h.services.Users.Register(c.Request.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			messageResponse(c, http.StatusConflict, "user already exists")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusCreated, UserIdResponse{UserId: resp})
}

// @Summary		User Login
// @Description	login and get token corresponding to user type
// @ID				userLogin
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		dtos.UserLoginInput	true	"User login info"
// @Success		200		{object}	authTokenResponse
// @Failure		400		{object}	response
// @Failure		404		{object}	response
// @Failure		500		{object}	response
// @Router			/auth/login [post]
func (h *Handler) userLogin(c *gin.Context) {
	var inp dtos.UserLoginInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")

		return
	}

	if err := inp.Validate(); err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	token, err := h.services.Users.Login(c.Request.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			messageResponse(c, http.StatusNotFound, "user not found")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, authTokenResponse{AuthToken: token})
}
