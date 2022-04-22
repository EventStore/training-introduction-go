package controllers

import (
	"net/http"
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/domain/writemodal/commands"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/labstack/echo/v4"
)

const (
	DateLayout = "2006-01-02"
)

type SlotsController struct {
	availableSlotsRepository readmodel.AvailableSlotsRepository
	patientSlotsRepository   readmodel.PatientSlotsRepository
	dispatcher               infrastructure.Dispatcher
}

func NewSlotsController(d infrastructure.Dispatcher, a readmodel.AvailableSlotsRepository, p readmodel.PatientSlotsRepository) *SlotsController {
	return &SlotsController{
		dispatcher:               d,
		availableSlotsRepository: a,
		patientSlotsRepository:   p,
	}
}

func (c *SlotsController) Register(e *echo.Echo) {
	e.GET("/slots/available/:date", c.AvailableHandler)
	e.GET("/slots/my-slots/:patientId", c.MySlotsHandler)
	e.POST("/slots/schedule", c.ScheduleHandler)
	e.POST("/slots/:slotId/book", c.BookHandler)
	e.POST("/slots/:slotId/cancel", c.CancelHandler)
}

func (c *SlotsController) AvailableHandler(ctx echo.Context) error {
	date, err := time.Parse(DateLayout, ctx.Param("date"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	availableSlots := c.availableSlotsRepository.GetSlotsAvailableOn(date)
	return ctx.JSON(http.StatusOK, availableSlots)
}

func (c *SlotsController) MySlotsHandler(ctx echo.Context) error {
	patientSlots := c.patientSlotsRepository.GetPatientSlots(ctx.Param("patientId"))
	return ctx.JSON(http.StatusOK, patientSlots)
}

func (c *SlotsController) ScheduleHandler(ctx echo.Context) error {
	req := ScheduleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := c.dispatcher.Dispatch(commands.NewScheduleCommand(req.SlotId, req.StartDateTime, req.Duration))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *SlotsController) BookHandler(ctx echo.Context) error {

	return ctx.NoContent(http.StatusOK)
}

func (c *SlotsController) CancelHandler(ctx echo.Context) error {

	return ctx.NoContent(http.StatusOK)
}

type ScheduleRequest struct {
	SlotId        string        `json:"slotId"`
	StartDateTime time.Time     `json:"startDateTime"`
	Duration      time.Duration `json:"duration"`
}

type BookRequest struct {
	PatientId string `json:"patientId"`
}

type CancelRequest struct {
	Reason string `json:"reason"`
}
