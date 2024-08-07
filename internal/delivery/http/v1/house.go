package v1

import (
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initHouseRoutes(api *gin.RouterGroup) {
	house := api.Group("/house")
	{
		authorized := house.Group("/", h.isAuthorized)
		{
			authorized.GET("/:id", h.getHouseById)
			authorized.POST("/:id/subscribe", h.postSubscribeToHouse)

			moderatorsOnly := authorized.Group("/", h.isModerator)
			{
				moderatorsOnly.POST("/create", h.createHouse)
			}
		}
	}
}

// @Summary		Get House By Id
// @Security		ClientsAuth
// @Security		ModeratorsAuth
// @Description	get all flats that are located at house
// @ID				getHouseById
// @Tags			house
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"house id"
// @Success		200	{object}	DataResponse[[]domain.Flat]
// @Failure		400	{object}	response
// @Failure		404	{object}	response
// @Failure		500	{object}	response
// @Router			/house/:id [get]
func (h *Handler) getHouseById(c *gin.Context) {
	houseId := c.Param("id")
	if houseId == "" {
		messageResponse(c, http.StatusBadRequest, "invalid house id")

		return
	}

	houseIdInt, err := strconv.Atoi(houseId)
	if err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid house id type")

		return
	}

	resp, err := h.services.Houses.GetById(c, houseIdInt)
	if err != nil {

		if errors.Is(err, domain.ErrHouseNotFound) {
			messageResponse(c, http.StatusNotFound, "house not found")

			return
		}

		messageResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, DataResponse[[]domain.Flat]{Data: resp})
}

// @Summary		Subscribe To House With Id
// @Security		ClientsAuth
// @Security		ModeratorsAuth
// @Description	subscribe user to house specifying his email in body
// @ID				postSubscribeToHouse
// @Tags			house
// @Accept			json
// @Produce		json
// @Param			id		path		string						true	"house id"
// @Param			input	body		dtos.HouseSubscribeInput	true	"User email to subscribe"
// @Success		200		{string}	string						"ok"
// @Failure		400		{object}	response
// @Failure		404		{object}	response
// @Failure		409		{object}	response
// @Failure		500		{object}	response
// @Router			/house/:id/subscribe [post]
func (h *Handler) postSubscribeToHouse(c *gin.Context) {
	houseId := c.Param("id")
	if houseId == "" {
		messageResponse(c, http.StatusBadRequest, "invalid house id")

		return
	}

	houseIdInt, err := strconv.Atoi(houseId)
	if err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid house id type")

		return
	}

	var email dtos.HouseSubscribeInput
	if err := c.BindJSON(&email); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")

		return
	}

	err = h.services.Houses.Subscribe(c, houseIdInt, email.Email)
	if err != nil {

		if errors.Is(err, domain.ErrHouseNotFound) {
			messageResponse(c, http.StatusNotFound, "house not found")

			return
		}

		if errors.Is(err, domain.ErrUserAlreadySubscribed) {
			messageResponse(c, http.StatusConflict, "user already subscribed")

			return
		}

		messageResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Create House
// @Security		ModeratorsAuth
// @Description	create house with address, year and (perhaps) developer
// @ID				createHouse
// @Tags			house
// @Accept			json
// @Produce		json
// @Param			input	body		dtos.HouseCreateInput	true	"House info"
// @Success		201		{object}	DataResponse[domain.House]
// @Failure		400		{object}	response
// @Failure		409		{object}	response
// @Failure		500		{object}	response
// @Router			/house/create [post]
func (h *Handler) createHouse(c *gin.Context) {
	var inp dtos.HouseCreateInput
	if err := c.BindJSON(&inp); err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	resp, err := h.services.Houses.Create(c.Request.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrHouseAlreadyExists) {
			messageResponse(c, http.StatusConflict, "house already exists")

			return
		}

		messageResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, DataResponse[domain.House]{Data: resp})
}
