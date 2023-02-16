package file

import (
	"context"
	"gotv/model"
	"gotv/resp"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
)

const BUCKET_NAME = "files"

type FileController struct {
	fileDao *fileDao
	context context.Context
	rc      *redis.Client
	mc      *minio.Client
}

func NewFileController(fileDao *fileDao, rc *redis.Client, mc *minio.Client) *FileController {
	return &FileController{
		fileDao: fileDao,
		context: context.Background(),
		rc:      rc,
		mc:      mc,
	}
}

func (f FileController) uplaod(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	print(file)
	if err != nil {
		resp.Fail(ctx, "出现异常: "+err.Error())
	}
	src, _ := file.Open()
	defer src.Close()

	objectName := strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + file.Filename

	buffer := make([]byte, 512)
	_, _ = src.Read(buffer)
	contentType := http.DetectContentType(buffer)
	src.Seek(0, 0)

	info, err := f.mc.PutObject(f.context, BUCKET_NAME, objectName, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Panicln(err)
	}
	var filePO model.File
	filePO.Filename = objectName
	filePO.Path = BUCKET_NAME + "/" + objectName
	filePO.Key = info.Key
	f.fileDao.SaveFile(filePO)
	resp.Success(ctx, info)
}

func (u *FileController) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	file := api.Group("/file")

	file.POST("/upload", u.uplaod)
}
