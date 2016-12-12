package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/file_api/models"
	"tuohai/internal/console"
	"tuohai/internal/file"
	"tuohai/internal/uuid"
)

func Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f, h, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "解析file文件失败"})
			return
		}
		cid := ctx.PostForm("cid")
		creator := ctx.PostForm("creator")
		if creator == "" || cid == "" {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "创建者或者cid不允许为空"})
			return
		}

		suffix := filepath.Ext(h.Filename)
		buf := &bytes.Buffer{}
		buf.ReadFrom(f)
		fmt.Println(buf.Len())
		l, _ := strconv.Atoi(ctx.PostForm("size"))
		if buf.Len() != l {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "上传意外终止"})
			return
		}

		finfo := &models.FileInfo{
			Id:       uuid.NewV4().String(),
			To:       cid,
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
		cid := ctx.Query("cid")
		tos := []string{cid}
		//获得to_id 查询文件
		info, err := models.GetFilesInfo(tos)
		if err != nil {
			ctx.String(200, "%s", "no")
			return
		}
		ctx.JSON(200, info)
	}
}
