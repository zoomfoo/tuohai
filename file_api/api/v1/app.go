package v1

import (
	// "errors"
	"bytes"
	"path/filepath"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/internal/file"
	"tuohai/internal/uuid"
	"tuohai/models"
)

func Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f, h, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.String(200, "%s", "no")
			return
		}
		to := ctx.PostForm("to")
		creator := ctx.PostForm("creator")

		suffix := filepath.Ext(h.Filename)

		buf := &bytes.Buffer{}
		buf.ReadFrom(f)

		finfo := &models.FileInfo{
			Id:       uuid.NewV4().String(),
			To:       to,
			Name:     h.Filename,
			Size:     len(buf.Bytes()),
			Type:     0,
			Ext:      suffix,
			Category: h.Header.Get("Content-Type"),
			Meta:     nil,
			Creator:  creator,
			Updated:  time.Now().Unix(),
			Created:  time.Now().Unix(),
		}

		if err := models.WriteFileToDB(finfo, file.UploadFile(suffix, buf)); err != nil {
			ctx.String(200, "%s", "no")
			console.StdLog.Error(err)
		} else {
			ctx.String(200, "%s", "ok")
		}
		return
	}
}

func Files() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//http获取用户相关的所有toid
		tos := []string{""}
		//获得to_id 查询文件
		info, err := models.GetFilesInfo(tos, models.FileTypeImage)
		if err != nil {
			ctx.String(200, "%s", "no")
			return
		}
		ctx.JSON(200, info)
	}
}
