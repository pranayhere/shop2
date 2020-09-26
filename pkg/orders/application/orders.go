package application

import (
	"github.com/pkg/errors"
	"log"
	"shop2/pkg/orders/domain/orders"
)

type OrdersService struct {
	ordersRepository orders.Repository
}

func NewOrderService(ordersRepository orders.Repository) OrdersService {
	return OrdersService{ordersRepository}
}

type PlaceOrderCommandAddress struct {
	Name     string
	Street   string
	City     string
	PostCode string
	Country  string
}

type PlaceOrderCommand struct {
	OrderID   orders.ID
	ProductID orders.ProductID

	Address PlaceOrderCommandAddress
}

func (s OrdersService) PlaceOrder(cmd PlaceOrderCommand) error {
	address, err := orders.NewAddress(
		cmd.Address.Name,
		cmd.Address.Street,
		cmd.Address.City,
		cmd.Address.PostCode,
		cmd.Address.Country,
	)
	if err != nil {
		return errors.Wrap(err, "invalid address")
	}

	product := orders.Product{}

	newOrder, err := orders.NewOrder(cmd.OrderID, product, address)
	if err != nil {
		return errors.Wrap(err, "cannot create order")
	}

	log.Printf("order %s placed", cmd.OrderID)
	log.Println(newOrder)
	return nil
}
