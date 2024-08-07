package v1

import (
	"bytes"
	"context"
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	mocks_service "github.com/dzhordano/avito-bootcamp2024/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_FlatCreate(t *testing.T) {
	type mockBehaviour func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput)

	tests := []struct {
		name               string
		inpBody            string
		inpFlat            dtos.FlatCreateInput
		mockBehaviour      mockBehaviour
		expectedStatusCode int
		expectedReqBody    string
	}{
		{
			name:    "OK",
			inpBody: `{"flat_number": 256, "house_id": 1, "price": 1000, "rooms": 3}`,
			inpFlat: dtos.FlatCreateInput{FlatNumber: 256, HouseId: 1, Price: 1000, Rooms: 3},
			mockBehaviour: func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput) {
				s.EXPECT().Create(context.Background(), inp).Return(domain.Flat{
					ID:         1,
					FlatNumber: 256,
					Price:      1000,
					Rooms:      3,
					Status:     domain.StatusCreated,
				}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedReqBody:    `{"data":{"ID":1,"FlatNumber":256,"Price":1000,"Rooms":3,"Status":"created"}}`,
		},
		{
			name:               "Empty body",
			inpBody:            "",
			inpFlat:            dtos.FlatCreateInput{},
			mockBehaviour:      func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedReqBody:    `{"message":"invalid input body"}`,
		},
		{
			name:    "Flat already exists",
			inpBody: `{"flat_number": 256, "house_id": 1, "price": 1000, "rooms": 3}`,
			inpFlat: dtos.FlatCreateInput{FlatNumber: 256, HouseId: 1, Price: 1000, Rooms: 3},
			mockBehaviour: func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput) {
				s.EXPECT().Create(context.Background(), inp).Return(domain.Flat{}, domain.ErrFlatAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
			expectedReqBody:    `{"message":"flat already exists"}`,
		},
		{
			name:    "House not found",
			inpBody: `{"flat_number": 256, "house_id": 100, "price": 1000, "rooms": 3}`,
			inpFlat: dtos.FlatCreateInput{FlatNumber: 256, HouseId: 100, Price: 1000, Rooms: 3},
			mockBehaviour: func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput) {
				s.EXPECT().Create(context.Background(), inp).Return(domain.Flat{}, domain.ErrHouseNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedReqBody:    `{"message":"house not found"}`,
		},
		{
			name:    "Internal server error",
			inpBody: `{"flat_number": 256, "house_id": 1, "price": 1000, "rooms": 3}`,
			inpFlat: dtos.FlatCreateInput{FlatNumber: 256, HouseId: 1, Price: 1000, Rooms: 3},
			mockBehaviour: func(s *mocks_service.MockFlats, inp dtos.FlatCreateInput) {
				s.EXPECT().Create(context.Background(), inp).Return(domain.Flat{}, errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReqBody:    `{"message":"internal server error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			flats := mocks_service.NewMockFlats(c)
			tt.mockBehaviour(flats, tt.inpFlat)

			services := &service.Services{
				Flats: flats,
			}
			handler := NewHandler(services, nil)

			r := gin.New()
			r.POST("/api/flat/create", handler.createFlat)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/api/flat/create", bytes.NewBufferString(tt.inpBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedReqBody, w.Body.String())
		})
	}
}
