package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"shop2/pkg/common/cmd"
	orders_app "shop2/pkg/orders/application"
	orders_infra "shop2/pkg/orders/infrastructure/orders"
	orders_public_http "shop2/pkg/orders/interfaces/public/http"
)

func main() {
	log.Println("Starting orders microservice")

	ctx := cmd.Context()

	r, closeFn := createOrdersMicroservice()
	defer closeFn()

	server := &http.Server{Addr: "127.0.0.1:8209", Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	log.Println("Closing orders microservice")

	if err := server.Close(); err != nil {
		panic(err)
	}
}

func createOrdersMicroservice() (router *chi.Mux, closeFn func()) {
	r := cmd.CreateRouter()

	ordersRepo := orders_infra.NewMemoryRepository()
	ordersService := orders_app.NewOrderService(ordersRepo)

	orders_public_http.AddRoutes(r, ordersService, ordersRepo)
	return r, func() {
		log.Println("closing application")
	}
}
