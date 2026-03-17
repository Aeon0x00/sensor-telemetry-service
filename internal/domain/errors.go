package domain

import "errors"

var (
	ErrMachineNotFound = errors.New("machine not found")
	ErrIncidentNotFound = errors.New("incident not found")
	ErrDuplicateEvent = errors.New("duplicate telemetry event")
	
	ErrInvalidMachine = errors.New("invalid machine data")
)
