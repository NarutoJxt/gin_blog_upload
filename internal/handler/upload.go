package handler

import (
	"gin_blog_upload/internal/model"
	"gin_blog_upload/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	fileService *service.FileService
}

func NewUploadHandler(fileService *service.FileService) *UploadHandler {
	return &UploadHandler{fileService: fileService}
}

func (h *UploadHandler) UploadFiles(c *gin.Context) {
	fileType := strings.TrimSpace(c.PostForm("type"))
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 7, Msg: "文件上传失败", Data: map[string]interface{}{}})
		return
	}

	items := make([]model.UploadItem, 0)
	for pos, fileList := range form.File {
		for _, fh := range fileList {
			fileURL, err := h.fileService.SaveUploadedFile(c.Request.Context(), fh, fileType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 7, Msg: "保存文件失败", Data: map[string]interface{}{}})
				return
			}
			items = append(items, model.UploadItem{
				URL: fileURL,
				POS: pos,
			})
		}
	}

	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Msg: "文件上传成功", Data: items})
}

