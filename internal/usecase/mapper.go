package usecase

import (
	"machine-telemetry-service/internal/domain"
	"math"
)

func mapEventToStateDTO (m *domain.Machine, ev domain.TelemetryEvent) domain.MachineStateDTO {
	return domain.MachineStateDTO{
		MachineId:  	m.Id,
		Code:       	m.Code,
		Workshop:   	m.Workshop,
		LastTempC:  	ev.TempC,
		LastVibro:  	ev.Vibro,
		LastRpm:    	ev.Rpm,
		LastPower:  	ev.PowerKw,
		IsOverload: 	ev.IsOverload,
	}
}


func calcShiftStats(events []domain.TelemetryEvent) (float64, float64, float64, float64, int, string) {
	if len(events) == 0 {
		return 0, 0, 0, 0, 0, "LOW"
	}

	var sumT, maxT float64
	var sumV, maxV float64
	var overCount int

	for _, ev := range events {
		sumT += ev.TempC
		sumV += ev.Vibro
		
		if ev.TempC > maxT {
			maxT = ev.TempC
		}
		if ev.Vibro > maxV {
			maxV = ev.Vibro
		}
		if ev.IsOverload {
			overCount++
		}
	}

	n := float64(len(events))
	avgT := math.Round((sumT/n)*100) / 100
	avgV := math.Round((sumV/n)*100) / 100

	risk := "LOW"
	if overCount >= 5 {
		risk = "HIGH"
	} else if overCount >= 1 {
		risk = "MEDIUM"
	}

	return avgT, maxT, avgV, maxV, overCount, risk
}
