package telemetry

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type RegisterMachineRequest struct {
	Code     string
	Workshop string
	EqType   string
	MaxTempC float64
	MaxVibro float64
}

type RegisterMachineResponse struct {
	MachineId string
}

type UpdateMachineLimitsRequest struct {
	MachineId string
	MaxTempC  float64
	MaxVibro  float64
}

type UpdateMachineLimitsResponse struct {
	Updated bool
}

type GetMachineStateRequest struct {
	MachineId string
}

type GetMachineStateResponse struct {
	MachineId  string
	Code       string
	Workshop   string
	LastTempC  float64
	LastVibro  float64
	LastRpm    int32
	IsOverload bool
}

type StartShiftReportRequest struct {
	MachineId string
}

type StartShiftReportResponse struct {
	ReportId string
	Status   string
}

type IncidentItem struct {
	Id        string
	MachineId string
	TempC     float64
	Vibro     float64
	Reason    string
}

type ListIncidentsRequest struct {
	MachineId string
}

type ListIncidentsResponse struct {
	Incidents []*IncidentItem
}

type ResolveIncidentRequest struct {
	IncidentId string
}

type ResolveIncidentResponse struct {
	Resolved bool
}

type TelemetryServiceServer interface {
	RegisterMachine(context.Context, *RegisterMachineRequest) (*RegisterMachineResponse, error)
	UpdateMachineLimits(context.Context, *UpdateMachineLimitsRequest) (*UpdateMachineLimitsResponse, error)
	GetMachineState(context.Context, *GetMachineStateRequest) (*GetMachineStateResponse, error)
	StartShiftReport(context.Context, *StartShiftReportRequest) (*StartShiftReportResponse, error)
	ListIncidents(context.Context, *ListIncidentsRequest) (*ListIncidentsResponse, error)
	ResolveIncident(context.Context, *ResolveIncidentRequest) (*ResolveIncidentResponse, error)
	mustEmbedUnimplementedTelemetryServiceServer()
}

type UnimplementedTelemetryServiceServer struct{}

func (UnimplementedTelemetryServiceServer) RegisterMachine(context.Context, *RegisterMachineRequest) (*RegisterMachineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMachine not implemented")
}
func (UnimplementedTelemetryServiceServer) UpdateMachineLimits(context.Context, *UpdateMachineLimitsRequest) (*UpdateMachineLimitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMachineLimits not implemented")
}
func (UnimplementedTelemetryServiceServer) GetMachineState(context.Context, *GetMachineStateRequest) (*GetMachineStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMachineState not implemented")
}
func (UnimplementedTelemetryServiceServer) StartShiftReport(context.Context, *StartShiftReportRequest) (*StartShiftReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartShiftReport not implemented")
}
func (UnimplementedTelemetryServiceServer) ListIncidents(context.Context, *ListIncidentsRequest) (*ListIncidentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListIncidents not implemented")
}
func (UnimplementedTelemetryServiceServer) ResolveIncident(context.Context, *ResolveIncidentRequest) (*ResolveIncidentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveIncident not implemented")
}
func (UnimplementedTelemetryServiceServer) mustEmbedUnimplementedTelemetryServiceServer() {}

func RegisterTelemetryServiceServer(s grpc.ServiceRegistrar, srv TelemetryServiceServer) {
	s.RegisterService(&TelemetryService_ServiceDesc, srv)
}

var TelemetryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "telemetry.v1.TelemetryService",
	HandlerType: (*TelemetryServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "proto/telemetry/v1/telemetry.proto",
}
