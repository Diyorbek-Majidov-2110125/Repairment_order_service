package storage

import (
	"context"
	"projects/Repairment_service/Repairment_order_service/genproto/order_service"
	"projects/Repairment_service/Repairment_order_service/models"

	"github.com/golang/protobuf/ptypes/empty"
)

type StorageI interface {
	CloseDB()
	Order() OrderRepoI
}

type OrderRepoI interface {
	Create(context.Context, *order_service.CreateOrderRequest) (*order_service.Order, error)
	GetById(context.Context, *order_service.OrderPrimaryKey) (*order_service.Order, error)
	GetList(context.Context, *order_service.GetListOrderRequest) (*order_service.GetListOrderResponse, error)
	Update(context.Context, *order_service.UpdateOrderRequest) (int64, error)
	UpdatePatch(context.Context, *models.UpdatePatchRequest) (int64, error)
	Delete(context.Context, *order_service.OrderPrimaryKey) (*empty.Empty,error)
}
