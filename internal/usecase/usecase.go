package usecase

import (
	"context"
	"fmt"
	"machine-telemetry-service/internal/domain"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase struct {
	Repo          	TelemetryRepository
	Cache         	TelemetryCache
	NotifyClient  	NotifyClient
	AlertProducer 	AlertProducer
	logger        	*zap.SugaredLogger
}

func NewUseCase (repo TelemetryRepository, cache TelemetryCache, notifyClient NotifyClient,alertProd AlertProducer, logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		Repo:          	repo,
		Cache:         	cache,
		NotifyClient:  	notifyClient,
		AlertProducer: 	alertProd,
		logger:        	logger,
	}
}
func (us *UseCase) RegisterMachine(ctx context.Context, code, workshop, eqType string, maxTempC, maxVibro float64) (string, error) {
	if code == "" || workshop == "" {
		return "", domain.ErrInvalidMachine
	}
	
	if maxTempC <= 0 {
		maxTempC = 85.0
	}
	if maxVibro <= 0 {
		maxVibro = 7.1
	}

	m := domain.Machine{
		Id:        	uuid.New().String(),
		Code:      	code,
		Workshop:  	workshop,
		EqType:    	eqType,
		MaxTempC:  	maxTempC,
		MaxVibro:  	maxVibro,
		CreatedAt: 	time.Now(),
	}

	if err := us.Repo.CreateMachine(ctx, m); err != nil {
		us.logger.Errorw(
			"failed to create machine",
			"code",  	code,
			"error", 	err,
		)
		return "", err
	}
	
	return  m.Id, nil
}

func (us *UseCase) UpdateMachineLimits(ctx context.Context, machineId string, maxTempC, maxVibro float64) error {
	if err := us.Repo.UpdateLimits(ctx, machineId, maxTempC, maxVibro); err != nil {
		return err
	}
	_ = us.Cache.DeleteMachineState(ctx, machineId)
	return nil
}


func (us *UseCase) ProcessKafkaTelemetry(ctx context.Context, dto domain.KafkaTelemetryDTO) error {
	fmt.Println("kafka dto:", dto.EventId, dto.MachineId)
	isNew, err := us.Cache.AcquireEventLock(ctx, dto.EventId)
	if err != nil {
		return err
	}
	if !isNew {
		us.logger.Infow("duplicate event skipped", "eventId", dto.EventId)
		return nil
	}

	machine,err := us.Repo.GetMachineById(ctx, dto.MachineId)
	if err != nil {
		_ = us.Cache.ReleaseEventLock(ctx, dto.EventId)
		return err
	}


	isOverload := dto.TempC >= machine.MaxTempC || dto.Vibro >= machine.MaxVibro

	ev := domain.TelemetryEvent{
		EventId:    	dto.EventId,
		MachineId:  	dto.MachineId,
		TempC:      	dto.TempC,
		Vibro:      	dto.Vibro,
		Rpm:        	dto.Rpm,
		PowerKw:    	dto.PowerKw,
		IsOverload: 	isOverload,
		CreatedAt:  	time.Now(),
	}

	if err := us.Repo.SaveTelemetryEvent(ctx, ev); err != nil {
		_ = us.Cache.ReleaseEventLock(ctx, dto.EventId)
		return err
	}

	if isOverload {
		streak, _ := us.Cache.IncrementOverloadStreak(ctx, machine.Id)
		if streak >= 3 {
			inc := domain.Incident{
				Id:         	uuid.New().String(),
				MachineId:  	machine.Id,
				EventId:    	ev.EventId,
				TempC:      	ev.TempC,
				Vibro:      	ev.Vibro,
				Reason:     	fmt.Sprintf("overload on %s: temp %.1fC, vibro %.2f mm/s", machine.Code, ev.TempC, ev.Vibro),
				IsResolved: 	false,
				CreatedAt:  	time.Now(),
			}
			_ = us.Repo.CreateIncident(ctx, inc)

			if us.NotifyClient != nil {
				us.NotifyClient.SendOverloadAlert(ctx, machine, inc)
			}
			if us.AlertProducer != nil {
				us.AlertProducer.PublishAlert(ctx, domain.KafkaAlertDTO{
					IncidentId: inc.Id,
					MachineId:  machine.Id,
					Code:       machine.Code,
					Workshop:   machine.Workshop,
					TempC:      inc.TempC,
					Vibro:      inc.Vibro,
					Reason:     inc.Reason,
				})
			}
		}
	} else {
		_ = us.Cache.ResetOverloadStreak(ctx, machine.Id)
	}

	stateDTO := mapEventToStateDTO(machine, ev)
	_ = us.Cache.SetMachineState(ctx, stateDTO)
	
	return nil
}

// func (us *UseCase) GetMachineStateFromDB(ctx context.Context, machineId string) (*domain.MachineStateDTO, error) {
// 	events, err := us.Repo.GetRecentEvents(ctx, machineId, 1)
// 	if err != nil || len(events) == 0 {
// 		// logger будет обрабатывать
// 		return nil, err
// 	}
// 	return &domain.MachineStateDTO{
// 		MachineId: machineId,
// 		LastTempC: events[0].TempC,
// 	}, nil
// }

func (us *UseCase) GetMachineState(ctx context.Context, machineId string) (*domain.MachineStateDTO, error) {
	cached, err := us.Cache.GetMachineState(ctx, machineId)
	if err == nil && cached != nil {
		return cached, nil
	}


	machine, err := us.Repo.GetMachineById(ctx, machineId)
	if err != nil {
		return nil, err
	}

	events, err := us.Repo.GetRecentEvents(ctx, machineId, 1)
	if err != nil || len(events) == 0 {
		return &domain.MachineStateDTO{
			MachineId: 	machine.Id,
			Code:      	machine.Code,
			Workshop:  	machine.Workshop,
		}, nil
	}

	stateDTO := mapEventToStateDTO(machine, events[0])
	_ = us.Cache.SetMachineState(ctx, stateDTO)

	return  &stateDTO, nil
}


func (us *UseCase) StartShiftReport(ctx context.Context, machineId string) (string, error) {
	reportId := uuid.New().String()
	
	rep := domain.ShiftReport{
		Id:        	reportId,
		MachineId: 	machineId,
		Status:    	"PENDING",
		RiskLevel: 	"LOW",
		CreatedAt: 	time.Now(),
	}
	if err := us.Repo.CreateShiftReport(ctx, rep); err != nil {
		return "", err
	}

	go us.buildShiftReport(reportId, machineId)
	return  reportId, nil
}

func (us *UseCase) buildShiftReport(reportId, machineId string) {
		bgCtx := context.Background()
		events, err := us.Repo.GetRecentEvents(bgCtx, machineId, 100)
		if err != nil {
			us.logger.Errorw(
				"failed to get recent events",
				"machineId", 	machineId,
				"reportId",  	reportId,
				"error",     	err,
			)
			return
		}
		avgT, maxT, avgV, maxV, overCnt, risk := calcShiftStats(events)
		
		rep := domain.ShiftReport{
			Id:            	reportId,
			MachineId:     	machineId,
			Status:        	"COMPLETED",
			AvgTempC:      	avgT,
			MaxTempC:      	maxT,
			AvgVibro:      	avgV,
			MaxVibro:      	maxV,
			OverloadCount: 	overCnt,
			RiskLevel:     	risk,
		}
		us.Repo.UpdateShiftReport(bgCtx, rep)
	}

func (us *UseCase) ListIncidents(ctx context.Context, machine_id string) ([]domain.Incident, error) {
	return us.Repo.ListActiveIncidents(ctx, machine_id)
}
func (us *UseCase) ResolveIncident(ctx context.Context, incidentId string) error {
	machineId, err := us.Repo.ResolveIncident(ctx, incidentId)
	if err != nil {
		return err
	}


	_ = us.Cache.DeleteMachineState(ctx, machineId)
	return nil
}
