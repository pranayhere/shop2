package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/satori/go.uuid"
	common_http "shop2/pkg/common/http"
	"shop2/pkg/orders/application"
	"shop2/pkg/orders/domain/orders"
)

func AddRoutes(router *chi.Mux, service application.OrdersService, repository orders.Repository) {
	resource := OrdersResource{service, repository}
	router.Post("/orders", resource.Create)
	router.Get("/orders/{id}/paid", resource.GetPaid)
}

type OrdersResource struct {
	service application.OrdersService

	repository orders.Repository
}

func (o OrdersResource) Create(w http.ResponseWriter, r *http.Request) {
	req := PostOrderRequest{}
	if err := render.Decode(r, &req); err != nil {
		_ = render.Render(w, r, common_http.ErrBadRequest(err))
		return
	}

	cmd := application.PlaceOrderCommand{
		OrderID: orders.ID(uuid.NewV1().String()),
		ProductID: req.ProductID,
		Address: application.PlaceOrderCommandAddress(req.Address),
	}

	if err := o.service.PlaceOrder(cmd); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, PostOrderResponse{
		OrderID: string(cmd.OrderID),
	})
}

func (o OrdersResource) GetPaid(w http.ResponseWriter, r *http.Request) {
	order, err := o.repository.ByID(orders.ID(chi.URLParam(r, "id")))
	if err != nil {
		_ = render.Render(w, r, common_http.ErrBadRequest(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	render.Respond(w, r, OrderPaidResponse{
		ID: string(order.ID()),
		IsPaid: order.Paid(),
	})
}

type PostOrderAddress struct {
	Name     string `json:"name"`
	Street   string `json:"street"`
	City     string `json:"city"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}

type PostOrderRequest struct {
	ProductID orders.ProductID `json:"product_id"`
	Address   PostOrderAddress `json:"address"`
}

type PostOrderResponse struct {
	OrderID string
}

type OrderPaidResponse struct {
	ID     string `json:"id"`
	IsPaid bool   `json:"is_paid"`
}
