package storage

import (
	"errors"
	"net/http"
	"time"

	"storage/app/utils"
	"storage/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	ID           string `json:"id"`
	ShortCode    string `json:"short_code"`
	PresignedURL string `json:"presigned_url"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (c *Controller) Upload(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		base.BadRequest(ctx, "file is required", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		base.InternalServerError(ctx, "cannot open file", nil)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	storage, err := c.svc.Upload(ctx.Request.Context(), fileHeader.Filename, contentType, fileHeader.Size, file)
	if err != nil {
		log.Errf("storage.upload.error: %v", err)
		if errors.Is(err, http.ErrMissingFile) {
			base.BadRequest(ctx, "file is required", nil)
			return
		}
		base.InternalServerError(ctx, "upload failed", nil)
		return
	}

	expiresIn := c.svc.Val.PresignExpireSeconds
	presignedURL, _, err := c.svc.Presign(ctx.Request.Context(), storage.ID, time.Duration(expiresIn)*time.Second)
	if err != nil {
		log.Errf("storage.upload.presign.error: %v", err)
		base.InternalServerError(ctx, "upload success but cannot create presign", nil)
		return
	}

	base.Success(ctx, UploadResponse{
		ID:           storage.ID.String(),
		ShortCode:    storage.ShortCode,
		PresignedURL: presignedURL,
		ExpiresIn:    expiresIn,
	})
}
