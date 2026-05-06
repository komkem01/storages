package storage

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"storage/app/utils"
	"storage/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PresignRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type PresignResponse struct {
	ID           string `json:"id"`
	Provider     string `json:"provider"`
	Path         string `json:"path"`
	URL          string `json:"url"`
	ShortCode    string `json:"short_code"`
	PresignedURL string `json:"presigned_url"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (c *Controller) Presign(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)

	var req PresignRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid id", nil)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid id", nil)
		return
	}

	expireSeconds := c.svc.Val.PresignExpireSeconds
	if raw := ctx.Query("expires_in"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v <= 0 {
			base.BadRequest(ctx, "invalid expires_in", nil)
			return
		}
		expireSeconds = v
	}

	url, storage, err := c.svc.Presign(ctx.Request.Context(), id, time.Duration(expireSeconds)*time.Second)
	if err != nil {
		log.Errf("storage.presign.error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			base.BadRequest(ctx, "storage not found", nil)
			return
		}
		base.InternalServerError(ctx, "presign failed", nil)
		return
	}

	base.Success(ctx, PresignResponse{
		ID:           storage.ID.String(),
		Provider:     storage.Provider,
		Path:         storage.Path,
		URL:          storage.URL,
		ShortCode:    storage.ShortCode,
		PresignedURL: url,
		ExpiresIn:    expireSeconds,
	})
}
