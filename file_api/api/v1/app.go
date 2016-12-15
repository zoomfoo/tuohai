package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/file_api/models"
	"tuohai/file_api/util"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/file"
	"tuohai/internal/uuid"
)

func Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.ParseMultipartForm(32 << 20)
		f, h, err := ctx.Request.FormFile("file")
		if err != nil {
			console.StdLog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "解析file文件失败"})
			return
		}

		cid := ctx.PostForm("cid")
		user := ctx.MustGet("user").(*auth.MainUser)
		creator := user.Uid

		if creator == "" || cid == "" {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "创建者或者cid不允许为空"})
			return
		}

		suffix := filepath.Ext(h.Filename)
		buf := &bytes.Buffer{}
		buf.ReadFrom(f)
		fmt.Println(buf.Len())
		now := time.Now().Unix()
		finfo := &models.FileInfo{
			Id:       uuid.NewV4().StringMd5(),
			To:       cid,
			Name:     h.Filename,
			Size:     buf.Len(),
			Ext:      suffix[1:],
			Category: h.Header.Get("Content-Type"),
			Meta:     nil,
			Creator:  creator,
			Updated:  now,
			Created:  now,
		}

		path := file.UploadFile(suffix, buf)
		width, height := 0, 0
		if util.IsImg(finfo.Category) {
			//获取图片宽高
			width, height = util.ImgDimension(buf)
			finfo.Meta = &models.Image{
				Id:         finfo.Id,
				ColorModel: "",
				Height:     height,
				Width:      width,
				Format:     "",
				Updated:    time.Now().Unix(),
				Created:    time.Now().Unix(),
			}
			finfo.Type = models.FileTypeImage
		}

		if err := models.WriteFileToDB(finfo, path); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1)
		} else {
			fmt.Println("Content-Type: ", h.Header.Get("Content-Type"))
			if util.IsImg(finfo.Category) {
				//如果是图片
				renderJSON(ctx, gin.H{
					"url":      path.P,
					"preview":  path.P,
					"type":     suffix[1:],
					"is_image": true,
					"width":    width,
					"height":   height,
					"owner":    creator,
					"time":     now,
				})
			} else {
				renderJSON(ctx, gin.H{
					"url":      path.P,
					"type":     suffix[1:],
					"is_image": false,
					"owner":    creator,
					"time":     now,
				})
			}

		}
		return
	}
}

func UploadAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f, h, err := ctx.Request.FormFile("file")
		if err != nil {
			console.StdLog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": struct{}{}, "message": "解析file文件失败"})
			return
		}
		buf := &bytes.Buffer{}
		buf.ReadFrom(f)
		suffix := filepath.Ext(h.Filename)
		res := file.UploadFile(suffix, buf)
		if res.E != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, "", 1, res.E.Error())
		} else {
			renderJSON(ctx, res.P)
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
			console.StdLog.Error(err)
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, info)
	}
}

func renderJSON(ctx *gin.Context, json interface{}, err_status ...interface{}) {
	switch len(err_status) {
	case 0:
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": json})
		break
	case 1:
		ctx.JSON(http.StatusOK, gin.H{"code": err_status[0], "data": json})
		break
	case 2:
		ctx.JSON(http.StatusOK, gin.H{"code": err_status[0], "msg": err_status[1], "data": json})
		break
	}
}
