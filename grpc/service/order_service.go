package service

import (
	"context"
	"projects/Repairment_service/Repairment_order_service/config"
	"projects/Repairment_service/Repairment_order_service/genproto/order_service"
	"projects/Repairment_service/Repairment_order_service/genproto/user_service"
	"projects/Repairment_service/Repairment_order_service/grpc/client"
	"projects/Repairment_service/Repairment_order_service/pkg/logger"
	"projects/Repairment_service/Repairment_order_service/storage"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	order_service.UnimplementedOrderServiceServer
}

func NewOrderService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) *orderService {
	return &orderService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: svcs,
	}
}

func (b *orderService) Create(ctx context.Context, req *order_service.CreateOrderRequest) (resp *order_service.Order, err error) {
	b.log.Info("---CreateOrder--->", logger.Any("req", req))

	user, err := b.services.UserService().GetById(ctx, &user_service.UserPrimaryKey{Id: req.UserId})
	
	if err != nil {
		b.log.Error("!!!CreateOrder--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = b.strg.Order().Create(ctx, req)
	if err != nil {
		b.log.Error("!!!CreateOrder--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp = &order_service.Order{
		Id:        resp.Id,
		UserId:    user.Id,
	}
	return resp, nil

}

func (b *orderService) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Order, err error) {
	b.log.Info("---GetOrderById--->", logger.Any("req", req))

	resp, err = b.strg.Order().GetById(ctx, req)
	if err != nil {
		b.log.Error("!!!GetOrderById--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (b *orderService) GetList(ctx context.Context, req *order_service.GetListOrderRequest) (resp *order_service.GetListOrderResponse, err error) {
	b.log.Info("---GetAllOrders--->", logger.Any("req", req))

	resp, err = b.strg.Order().GetList(ctx, req)
	if err != nil {
		b.log.Error("!!!GetAllOrders--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (b *orderService) Update(ctx context.Context, req *order_service.UpdateOrderRequest) (resp *order_service.Order, err error) {
	b.log.Info("---UpdateOrder--->", logger.Any("req", req))

	_, err = b.strg.Order().Update(ctx, req)
	if err != nil {
		b.log.Error("!!!UpdateOrder--->", logger.Error(err))
		return resp, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (b *orderService) Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (*empty.Empty, error) {
	b.log.Info("---DeleteOrder--->", logger.Any("req", req))

	 resp, err := b.strg.Order().Delete(ctx, req)
	if err != nil {
		b.log.Error("!!!DeleteOrder--->", logger.Error(err))
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}



func (b *orderService) GetUserInfoById(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *user_service.User, err error) {
	order := &order_service.Order{}
	order, err = b.strg.Order().GetById(ctx, req)
	if err != nil {
		b.log.Error("!!!GetUserInfoById--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := b.services.UserService().GetById(ctx, &user_service.UserPrimaryKey{Id: order.UserId})

	if err != nil {
		b.log.Error("!!!GetUserInfoById--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return user, nil

}
