package grpcclient

import (
	"machine-telemetry-service/gen/notify"
	"machine-telemetry-service/internal/domain"
)

func mapAlertToProto (machine *domain.Machine, inc domain.Incident) *notify.SendAlertRequest {
	return &notify.SendAlertRequest{
		MachineId: 	machine.Id,
		Code:      	machine.Code,
		Workshop:  	machine.Workshop,
		Reason:    	inc.Reason,
		TempC:     	inc.TempC,
		Vibro:     	inc.Vibro,
	}
}
