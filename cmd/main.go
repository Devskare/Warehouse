package main

import (
	"log/slog"
	"net"
	"time"
	"warehouse/config"
	warehousev1 "warehouse/gen/warehouse/v1"
	"warehouse/logger"
	"warehouse/modules/Whouse/Whousegrpc"
	"warehouse/modules/Whouse/repository"
	"warehouse/modules/Whouse/service"
	"warehouse/modules/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	time.Sleep(2 * time.Second)

	appConf := config.MustLoadConfig(".env")
	log := logger.Initlogger(appConf.LogLevel, appConf.Production)

	sqlDB, err := db.NewSqlDB(log, &appConf.DB)
	if err != nil {
		log.Error("failed to connect sql data base", slog.Any("error", err))
		panic(err)
	}

	repo := repository.NewWHouseRepository(sqlDB)
	ProductService := service.NewProductService(repo, log)
	StorageService := service.NewStorageService(repo, log)
	s := initRPC(ProductService, StorageService)
	lis, err := net.Listen("tcp", appConf.GrpcServerPort)
	if err != nil {
		log.Error("failed to listen", slog.Any("error", err))
		panic(err)
	}
	log.Info("grpc server listening on", slog.Any("address", lis.Addr().String()))
	if err = s.Serve(lis); err != nil {
		log.Error("failed to serve", slog.Any("error", err))
		panic(err)
	}
}

func initRPC(productService *service.ProductService, storageService *service.StorageService) *grpc.Server {
	s := grpc.NewServer()

	warehousev1.RegisterWarehouseServer(
		s, Whousegrpc.NewWarehouseServer(productService, storageService),
	)
	reflection.Register(s)
	return s
}
