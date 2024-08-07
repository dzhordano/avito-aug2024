package v1

import (
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initFlatRoutes(api *gin.RouterGroup) {
	flats := api.Group("/flat")
	{
		authorized := flats.Group("/", h.isAuthorized)
		{
			authorized.POST("/create", h.createFlat)

			moderatorsOnly := authorized.Group("/", h.isModerator)
			{
				moderatorsOnly.POST("/update", h.updateFlat)
			}
		}
	}
}

// @Summary		Create flat
// @Security		ClientsAuth
// @Security		ModeratorsAuth
// @Description	create flat with flatNumber, price, rooms and house id it belongs to
// @ID				createFlat
// @Tags			flat
// @Accept			json
// @Produce		json
// @Param			input	body		dtos.FlatCreateInput	true	"Flat info"
// @Success		201		{object}	DataResponse[domain.Flat]
// @Failure		400		{object}	response
// @Failure		409		{object}	response
// @Failure		500		{object}	response
// @Router			/flat/create [post]
func (h *Handler) createFlat(c *gin.Context) {
	var inp dtos.FlatCreateInput
	if err := c.BindJSON(&inp); err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	resp, err := h.services.Flats.Create(c.Request.Context(), inp)
	if err != nil {

		if errors.Is(err, domain.ErrFlatAlreadyExists) {
			messageResponse(c, http.StatusConflict, "flat already exists")

			return
		}

		if errors.Is(err, domain.ErrHouseNotFound) {
			messageResponse(c, http.StatusNotFound, "house not found")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusCreated, DataResponse[domain.Flat]{Data: resp})
}

// @Summary		Update flat
// @Security		ModeratorsAuth
// @Description	update flat status
// @ID				updateFlat
// @Tags			flat
// @Accept			json
// @Produce		json
// @Param			input	body		dtos.FlatUpdateInput	true	"Flat info"
// @Success		200		{object}	DataResponse[domain.Flat]
// @Failure		400		{object}	response
// @Failure		404		{object}	response
// @Failure		500		{object}	response
// @Router			/flat/update [post]
func (h *Handler) updateFlat(c *gin.Context) {
	var inp dtos.FlatUpdateInput
	if err := c.BindJSON(&inp); err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := inp.Validate(); err != nil {
		messageResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := h.services.Flats.Update(c.Request.Context(), inp.FlatId, inp.Status)
	if err != nil {
		if errors.Is(err, domain.ErrFlatNotFound) {
			messageResponse(c, http.StatusNotFound, "flat not found")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, DataResponse[domain.Flat]{Data: resp})
}
