package main

import (
	"context"
	"gin_blog_upload/internal/config"
	"gin_blog_upload/internal/handler"
	"gin_blog_upload/internal/middleware"
	"gin_blog_upload/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.New()
	router.Use(gin.Logger(), middleware.RecoverMiddleware())
	router.MaxMultipartMemory = 16 << 20

	fileService, err := service.NewFileService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
		cfg.MinioPublicBase,
	)
	if err != nil {
		panic(err)
	}
	if err := fileService.EnsureBucket(context.Background()); err != nil {
		panic(err)
	}
	uploadHandler := handler.NewUploadHandler(fileService)

	// 上传接口（文件实际存储在 MinIO）
	router.POST("/api/v1/files", uploadHandler.UploadFiles)

	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		panic(err)
	}
}

