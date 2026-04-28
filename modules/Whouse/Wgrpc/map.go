package Wgrpc

import (
	warehousev1 "warehouse/gen/warehouse/v1"
	"warehouse/modules/Whouse/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapProductModelToProto(p *models.ProductModel) *warehousev1.Product {
	if p == nil {
		return nil
	}

	Proto := &warehousev1.Product{
		Article: int64(p.Article),
		Weight:  p.Weight,
		Name:    p.ProductName,
	}

	if p.StorageID != nil {
		Proto.StorageId = int64(*p.StorageID)
	}

	if p.DeliveryDate != nil {
		Proto.DeliveryDate = timestamppb.New(*p.DeliveryDate)
	}

	if p.ExpireDate != nil {
		Proto.ExpireDate = timestamppb.New(*p.ExpireDate)
	}

	return Proto

}

func mapProductProtoToModel(p *warehousev1.Product) *models.ProductModel {
	if p == nil {
		return nil
	}

	model := &models.ProductModel{
		Article:     int(p.GetArticle()),
		ProductName: p.GetName(),
		Weight:      p.GetWeight(),
	}

	if p.GetStorageId() != 0 {
		storageID := int(p.GetStorageId())
		model.StorageID = &storageID
	}

	if p.GetDeliveryDate() != nil {
		t := p.GetDeliveryDate().AsTime()
		model.DeliveryDate = &t
	}

	if p.GetExpireDate() != nil {
		t := p.GetExpireDate().AsTime()
		model.ExpireDate = &t
	}

	return model
}
