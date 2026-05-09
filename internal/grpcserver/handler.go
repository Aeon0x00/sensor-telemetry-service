package grpcserver

import (
	"context"
	"machine-telemetry-service/gen/telemetry"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (s *Server) RegisterMachine(ctx context.Context, req *telemetry.RegisterMachineRequest) (*telemetry.RegisterMachineResponse, error) {
	machine_id, err := s.us.RegisterMachine(ctx, req.Code, req.Workshop, req.EqType, req.MaxTempC, req.MaxVibro)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register machine: %v", err)
	}
	return &telemetry.RegisterMachineResponse{
		MachineId: machine_id,
	}, nil
}
func (s *Server) UpdateMachineLimits(ctx context.Context, req *telemetry.UpdateMachineLimitsRequest) (*telemetry.UpdateMachineLimitsResponse, error) {
	if err := s.us.UpdateMachineLimits(ctx, req.MachineId, req.MaxTempC, req.MaxVibro); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update limits: %v", err)
	}

	return &telemetry.UpdateMachineLimitsResponse{
		Updated: true,
	}, nil
}
func (s *Server) GetMachineState(ctx context.Context, req *telemetry.GetMachineStateRequest) (*telemetry.GetMachineStateResponse, error) {
	state, err := s.us.GetMachineState(ctx, req.MachineId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get machine state: %v", err)
	}


	return &telemetry.GetMachineStateResponse{
		MachineId:  state.MachineId,
		Code:       state.Code,
		Workshop:   state.Workshop,
		LastTempC:  state.LastTempC,
		LastVibro:  state.LastVibro,
		LastRpm:    int32(state.LastRpm),
		IsOverload: state.IsOverload,
	}, nil
}


func (s *Server) StartShiftReport(ctx context.Context, req *telemetry.StartShiftReportRequest) (*telemetry.StartShiftReportResponse, error) {
	machineId := req.MachineId

	report_id, err := s.us.StartShiftReport(ctx, machineId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start shift report: %v", err)
	}
	return &telemetry.StartShiftReportResponse{
		ReportId: report_id,
		Status:   "REPORT_STATUS_PENDING",
	}, nil
}

func (s *Server) ListIncidents(ctx context.Context, req *telemetry.ListIncidentsRequest) (*telemetry.ListIncidentsResponse, error) {
	list, err := s.us.ListIncidents(ctx, req.MachineId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list incidents: %v", err)
	}

	items := []*telemetry.IncidentItem{}
	for _, inc := range list{
		items = append(items, &telemetry.IncidentItem{
			Id:        inc.Id,
			MachineId: inc.MachineId,
			TempC:     inc.TempC,
			Vibro:     inc.Vibro,
			Reason:    inc.Reason,
		})
	}
	return  &telemetry.ListIncidentsResponse{
		Incidents: items,
	}, nil
}
func (s *Server) ResolveIncident(ctx context.Context, req *telemetry.ResolveIncidentRequest) (*telemetry.ResolveIncidentResponse, error) {
	if err := s.us.ResolveIncident(ctx, req.IncidentId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to resolve incident: %v", err)
	}

	return &telemetry.ResolveIncidentResponse{
		Resolved: true,
	}, nil
}
