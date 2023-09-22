package httpcontroller

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/service"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestController_GetOrder(t *testing.T) {

	tests := []struct {
		name              string
		setupServiceMock  func(sm *ServiceMock, orderId string)
		setupTemplateMock func(tem *TemplateExecutorMock)
		reqOrderUID       string
		expectedCode      int
	}{
		{
			name: "OK",
			setupServiceMock: func(sm *ServiceMock, orderId string) {
				sm.GetOrderMock.Expect(context.Background(), orderId).Return(models.Order{OrderUID: orderId}, nil)
			},
			setupTemplateMock: func(tem *TemplateExecutorMock) {
				tem.ExecuteTemplateMock.Expect(&bytes.Buffer{}, "order.html", models.Order{OrderUID: "existedOrder"}).Return(nil)
			},
			reqOrderUID:  "existedOrder",
			expectedCode: 200,
		},
		{
			name: "order not found",
			setupServiceMock: func(sm *ServiceMock, orderId string) {
				sm.GetOrderMock.Expect(context.Background(), orderId).Return(models.Order{}, service.ErrOrderNotFound)
			},
			setupTemplateMock: func(tem *TemplateExecutorMock) {
				tem.ExecuteTemplateMock.Expect(&bytes.Buffer{}, "not_found.html", "job").Return(nil)
			},
			reqOrderUID:  "job",
			expectedCode: 404,
		},
		{
			name: "service unexpected error",
			setupServiceMock: func(sm *ServiceMock, orderId string) {
				sm.GetOrderMock.Expect(context.Background(), orderId).Return(models.Order{}, errors.New("unexpected error"))
			},
			setupTemplateMock: func(tem *TemplateExecutorMock) {
			},
			reqOrderUID:  "fatal",
			expectedCode: 500,
		},
		{
			name: "error rendering",
			setupServiceMock: func(sm *ServiceMock, orderId string) {
				sm.GetOrderMock.Expect(context.Background(), orderId).Return(models.Order{OrderUID: orderId}, nil)
			},
			setupTemplateMock: func(tem *TemplateExecutorMock) {
				tem.ExecuteTemplateMock.Expect(&bytes.Buffer{}, "order.html", models.Order{OrderUID: "unimaginable"}).Return(errors.New("cannot parse"))
			},
			reqOrderUID:  "unimaginable",
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewServiceMock(t)
			tem := NewTemplateExecutorMock(t)

			c := Controller{
				service:  sm,
				template: tem,
			}

			tt.setupServiceMock(sm, tt.reqOrderUID)
			tt.setupTemplateMock(tem)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/orders/"+tt.reqOrderUID, nil)
			p := httprouter.Params{
				httprouter.Param{
					Key:   "OrderUID",
					Value: tt.reqOrderUID,
				},
			}

			c.GetOrder(w, r, p)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
