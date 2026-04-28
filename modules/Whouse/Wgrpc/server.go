package Wgrpc

import (
	"context"
	warehousev1 "warehouse/gen/warehouse/v1"
	"warehouse/modules/Whouse/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	productModel := mapProductProtoToModel(req.GetProduct())
	if productModel == nil {
		return nil, status.Error(codes.Internal, "failed to map product")
	}
	err := s.productService.ProductADD(ctx, *productModel)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) UpdateProduct(ctx context.Context, req *warehousev1.UpdateProductRequest) (*warehousev1.Empty, error) {
	if req.GetProduct() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "product is required")
	}

	productModel := mapProductProtoToModel(req.GetProduct())
	if productModel == nil {
		return nil, status.Error(codes.Internal, "failed to map product")
	}
	err := s.productService.ProductUpdate(ctx, *productModel)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &warehousev1.Empty{}, nil

}

func (s *WarehouseRPCServer) DeleteProduct(ctx context.Context, req *warehousev1.DeleteProductRequest) (*warehousev1.Empty, error) {
	err := s.productService.ProductDelete(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) ExpireProduct(ctx context.Context, req *warehousev1.ExpireProductRequest) (*warehousev1.Empty, error) {
	err := s.productService.ProductExpire(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &warehousev1.Empty{}, nil

}

func (s *WarehouseRPCServer) GetProduct(ctx context.Context, req *warehousev1.GetProductRequest) (*warehousev1.Product, error) {
	product, err := s.productService.GetProductFromProducts(ctx, int(req.GetArticle()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if product == nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return mapProductModelToProto(product), nil

}

func (s *WarehouseRPCServer) ListProducts(ctx context.Context, req *warehousev1.Empty) (*warehousev1.ListProductsResponse, error) {
	products, err := s.productService.ProductsListFromProducts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &warehousev1.Empty{}, nil
}

func (s *WarehouseRPCServer) ListStorages(ctx context.Context, req *warehousev1.Empty) (*warehousev1.ListStoragesResponse, error) {
	storages, err := s.storageService.ListStorages(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &warehousev1.ListStoragesResponse{
		Storages: make([]*warehousev1.Storage, 0, len(storages)),
	}
	for _, storage := range storages {
		resp.Storages = append(resp.Storages, &warehousev1.Storage{
			Id:            int64(storage.ID),
			MaxWeight:     storage.MaxWeight,
			CurrentWeight: storage.CurrentWeight,
		})
	}

	return resp, nil
}
