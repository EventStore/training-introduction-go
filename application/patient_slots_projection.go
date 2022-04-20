package application

import (
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
)

type PatientSlotsProjection struct {
	projections.ProjectionBase

	repository readmodel.PatientSlotsRepository
}

func NewPatientSlotsProjection(r readmodel.PatientSlotsRepository) *PatientSlotsProjection {
	p := projections.NewProjection()

	return &PatientSlotsProjection{p, r}
}
