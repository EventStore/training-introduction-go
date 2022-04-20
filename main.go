package main

import (
	"context"
	"net/http"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/application"
	"github.com/EventStore/training-introduction-go/controllers"
	"github.com/EventStore/training-introduction-go/domain/writemodal"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/inmemory"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := createESDBClient()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	f := infrastructure.NewEventFactory(events.Booked{}, events.Cancelled{}, events.Scheduled{})
	eventStore := infrastructure.NewEsEventStore(db, f)
	aggregateStore := infrastructure.NewEsAggregateStore(eventStore)
	asr := inmemory.NewInMemoryAvailableSlotsRepository()
	psr := inmemory.NewInMemoryPatientSlotsRepository()
	subManager := projections.NewSubscriptionManager(db, f,
		projections.NewProjector(application.NewAvailableSlotsProjection(asr)))
	subManager.Start(context.Background())

	c := infrastructure.NewCommandHandlerMap(writemodal.NewSlotHandlers(aggregateStore))
	d := infrastructure.NewDispatcher(c)

	s := controllers.NewSlotsController(d, asr, psr)
	s.Register(e)

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":5001"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Training!")
}

func createESDBClient() (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString("esdb://localhost:2113?tls=false")
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}
