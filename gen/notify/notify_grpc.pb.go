package notify

import (
	context "context"
	grpc "google.golang.org/grpc"
)

type SendAlertRequest struct {
	MachineId string
	Code      string
	Workshop  string
	Reason    string
	TempC     float64
	Vibro     float64
}

type SendAlertResponse struct {
	Delivered bool
}

type NotifyServiceClient interface {
	SendAlert(ctx context.Context, in *SendAlertRequest, opts ...grpc.CallOption) (*SendAlertResponse, error)
}

type notifyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotifyServiceClient(cc grpc.ClientConnInterface) NotifyServiceClient {
	return &notifyServiceClient{cc}
}

func (c *notifyServiceClient) SendAlert(ctx context.Context, in *SendAlertRequest, opts ...grpc.CallOption) (*SendAlertResponse, error) {
	out := new(SendAlertResponse)
	err := c.cc.Invoke(ctx, "/notify.v1.NotifyService/SendAlert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
