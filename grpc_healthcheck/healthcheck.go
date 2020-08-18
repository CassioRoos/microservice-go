package grpc_healthcheck

import (
	"context"
	"github.com/CassioRoos/grpc_currency/protos/healthcheck"
	"github.com/hashicorp/go-hclog"
	"os"
	"time"
)

type grpcHealthCheck struct {
	log hclog.Logger
	h  healthcheck.HealthCheckClient
}

func NewGrpcHealthCheck(log hclog.Logger, h healthcheck.HealthCheckClient) *grpcHealthCheck{
	return &grpcHealthCheck{log: log, h: h}
}

func (h *grpcHealthCheck) HealthCheck(times int) bool {
	for i := 0; i <= times; i ++{
		 resp, err := h.h.Check(context.Background(), &healthcheck.HealthCheckParam{})
		 if err != nil && i == times{
		 	h.log.Error("GRPC is unhealty and wont connect", "ERROR", err)
		 	os.Exit(1)
		 }
		 if err == nil{
		 	h.log.Info("GRPC is healthy", "Message", resp.Message)
		 	return true
		 }
		 time.Sleep(1 * time.Second)
	}
	return false
}