package usecase

import (
	"context"
	"machine-telemetry-service/internal/domain"

)

type NotifyClient interface {
	SendOverloadAlert(ctx context.Context, machine *domain.Machine, inc domain.Incident) error
}

type AlertProducer interface {
	PublishAlert(ctx context.Context, alert domain.KafkaAlertDTO) error
}

type TelemetryRepository interface {
	CreateMachine(ctx context.Context, m domain.Machine) error
	UpdateLimits(ctx context.Context, machineId string, maxTempC, maxVibro float64) error
	GetMachineById(ctx context.Context, machineId string) (*domain.Machine, error)
	SaveTelemetryEvent(ctx context.Context, ev domain.TelemetryEvent) error
	
	GetRecentEvents(ctx context.Context, machineId string, limit int) ([]domain.TelemetryEvent, error)
	CreateShiftReport(ctx context.Context, rep domain.ShiftReport) error
	UpdateShiftReport(ctx context.Context, rep domain.ShiftReport) error
	CreateIncident(ctx context.Context, inc domain.Incident) error
	ListActiveIncidents(ctx context.Context, machine_id string) ([]domain.Incident, error)
	ResolveIncident(ctx context.Context, incidentId string) (string, error)
}

type TelemetryCache interface {
	AcquireEventLock(ctx context.Context, eventId string) (bool, error)
	ReleaseEventLock(ctx context.Context, eventId string) error
	IncrementOverloadStreak(ctx context.Context, machineId string) (int64, error)
	ResetOverloadStreak(ctx context.Context, machineId string) error

	GetMachineState(ctx context.Context, machineId string) (*domain.MachineStateDTO, error)
	SetMachineState(ctx context.Context, state domain.MachineStateDTO) error
	DeleteMachineState(ctx context.Context, machineId string) error
}
