package grpc

import (
	"projects/Repairment_service/Repairment_order_service/config"
	"projects/Repairment_service/Repairment_order_service/genproto/order_service"
	"projects/Repairment_service/Repairment_order_service/grpc/client"
	"projects/Repairment_service/Repairment_order_service/grpc/service"
	"projects/Repairment_service/Repairment_order_service/pkg/logger"
	"projects/Repairment_service/Repairment_order_service/storage"

	"google.golang.org/grpc"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	order_service.RegisterOrderServiceServer(grpcServer, service.NewOrderService(cfg, log, strg, svcs))

	return grpcServer
}
