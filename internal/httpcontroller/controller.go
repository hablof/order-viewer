package httpcontroller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hablof/order-viewer/internal/app/service"
	"github.com/hablof/order-viewer/internal/models"
	"github.com/julienschmidt/httprouter"
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
	order, err := c.service.GetOrder(r.Context(), OrderUID)
	switch {
	case err == service.ErrOrderNotFound:
		c.renderNotFound(w, r, OrderUID)
		return

	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	b := bytes.Buffer{}
	if err := c.template.ExecuteTemplate(&b, "order.html", order); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")

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

	w.Write(b.Bytes())
}
