package httpcontroller

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"

	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/service"
)

type Service interface {
	GetOrder(ctx context.Context, OrderUID string) (models.Order, error)
}

type TemplateExecutor interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

type Controller struct {
	service  Service
	template TemplateExecutor
}

func GetRouter(s Service, t TemplateExecutor) http.Handler {
	c := Controller{
		service:  s,
		template: t,
	}

	router := httprouter.New()
	router.GET("/orders/:OrderUID", c.GetOrder)

	return router
}

// основной http handle
// отрисовывает информацию о заказе, либо страницу 404
func (c *Controller) GetOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	OrderUID := p.ByName("OrderUID")

	log := log.Logger.With().Str("func", "Controller.GetOrder").Caller().Str("requestedOrderUID", OrderUID).Logger()

	log.Debug().Msg("request to get order")

	order, err := c.service.GetOrder(r.Context(), OrderUID)
	switch {
	case err == service.ErrOrderNotFound:
		log.Debug().Msg("order not found")
		c.renderNotFound(w, r, OrderUID)

		return

	case err != nil:
		log.Error().Err(err).Msg("unexpected error")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))

		return
	}

	b := bytes.Buffer{}
	if err := c.template.ExecuteTemplate(&b, "order.html", order); err != nil {
		log.Error().Err(err).Msg("template execution unepected error")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}

func (c *Controller) renderNotFound(w http.ResponseWriter, r *http.Request, OrderUID string) {
	b := bytes.Buffer{}
	if err := c.template.ExecuteTemplate(&b, "not_found.html", OrderUID); err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(b.Bytes())
}
