package Whousegrpc

import (
	"context"
	warehousev1 "warehouse/gen/warehouse/v1"
	"warehouse/modules/Whouse/models"
	"warehouse/modules/Whouse/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WarehouseRPCServer struct {
	warehousev1.UnimplementedWarehouseServer
	productService *service.ProductService
	storageService *service.StorageService
}

func NewWarehouseServer(
	ps *service.ProductService,
	ss *service.StorageService,
) *WarehouseRPCServer {
	return &WarehouseRPCServer{
		productService: ps,
		storageService: ss,
	}
}

func (s *WarehouseRPCServer) AddProduct(ctx context.Context, req *warehousev1.AddProductRequest) (*warehousev1.Empty, error) {
	if req.GetProduct() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "product is required")
	}

	p := req.GetProduct()

	storageID := int(p.GetStorageId())

	productModel := models.ProductModel{
		Article:     int(p.GetArticle()),
		ProductName: p.GetName(),
		StorageID:   &storageID,
		Weight:      p.GetWeight(),
	}

	err := s.productService.ProductADD(ctx, productModel)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) UpdateProduct(ctx context.Context, req *warehousev1.UpdateProductRequest) (*warehousev1.Empty, error) {
	if req.GetProduct() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "product is required")
	}
	p := req.GetProduct()

	storageID := int(p.GetStorageId())

	productModel := models.ProductModel{
		Article:     int(p.GetArticle()),
		ProductName: p.GetName(),
		StorageID:   &storageID,
		Weight:      p.GetWeight(),
	}

	err := s.productService.ProductUpdate(ctx, productModel)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &warehousev1.Empty{}, nil

}

func (s *WarehouseRPCServer) DeleteProduct(ctx context.Context, req *warehousev1.DeleteProductRequest) (*warehousev1.Empty, error) {
	err := s.productService.ProductDelete(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, err
	}

	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) ExpireProduct(ctx context.Context, req *warehousev1.ExpireProductRequest) (*warehousev1.Empty, error) {
	err := s.productService.ProductExpire(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, err
	}
	return &warehousev1.Empty{}, nil

}

func (s *WarehouseRPCServer) GetProduct(ctx context.Context, req *warehousev1.GetProductRequest) (*warehousev1.Product, error) {
	product, err := s.productService.GetProductFromProducts(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, status.Error(codes.NotFound, "product is nil")
	}
	return mapProductModelToProto(product), nil

}

func (s *WarehouseRPCServer) ListProducts(ctx context.Context, req *warehousev1.Empty) (*warehousev1.ListProductsResponse, error) {
	products, err := s.productService.ProductsListFromProducts(ctx)
	if err != nil {
		return nil, err
	}

	resp := &warehousev1.ListProductsResponse{
		Products: make([]*warehousev1.Product, 0, len(products)),
	}
	for _, p := range products {
		protoProduct := mapProductModelToProto(&p)
		if protoProduct != nil {
			resp.Products = append(resp.Products, protoProduct)
		}
	}

	return resp, nil
}

func (s *WarehouseRPCServer) AddStorage(ctx context.Context, req *warehousev1.AddStorageRequest) (*warehousev1.Empty, error) {
	err := s.storageService.StorageADD(ctx, req.GetMaxWeight())
	if err != nil {
		return nil, err
	}
	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) ListStorages(ctx context.Context, req *warehousev1.Empty) (*warehousev1.ListStoragesResponse, error) {
	storages, err := s.storageService.ListStorages(ctx)
	if err != nil {
		return nil, err
	}
	resp := &warehousev1.ListStoragesResponse{
		Storages: make([]*warehousev1.Storage, 0, len(storages)),
	}
	for _, storage := range storages {
		resp.Storages = append(resp.Storages, &warehousev1.Storage{
			Id:        int64(storage.ID),
			MaxWeight: storage.MaxWeight,
		})
	}

	return resp, nil
}

func mapProductModelToProto(p *models.ProductModel) *warehousev1.Product {
	if p == nil {
		return nil
	}

	res := &warehousev1.Product{
		Article: int64(p.Article),
		Weight:  p.Weight,
		Name:    p.ProductName,
	}

	if p.StorageID != nil {
		res.StorageId = int64(*p.StorageID)
	}

	if p.DeliveryDate != nil {
		res.DeliveryDate = timestamppb.New(*p.DeliveryDate)
	}

	if p.ExpireDate != nil {
		res.ExpireDate = timestamppb.New(*p.ExpireDate)
	}

	return res

}
