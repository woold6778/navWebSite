package upload

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"nav-web-site/mydb"
	"nav-web-site/util"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ImgReturnData struct {
	Hash    string `json:"hash"`
	ImgPath string `json:"img_path"`
}

// UploadImage 处理图片上传
// @Summary 上传图片
// @Description 上传图片文件
// @Tags upload
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "图片文件"
// @Success 200 {object} util.APIResponse{code=int,message=string,data=ImgReturnData}
// @Failure 400 {object} util.APIResponse{code=int,message=string,data=object}
// @Router /upload/image [post]
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "获取文件失败", Data: "null"})
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "不支持的文件格式", Data: "null"})
		return
	}

	// 打开文件并检查是否为合法图片
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "读取文件内容失败", Data: "null"})
		return
	}
	defer fileContent.Close()

	// 尝试解码图片
	_, _, err = image.Decode(fileContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse{Code: http.StatusBadRequest, Message: "文件不是合法的图片", Data: "null"})
		return
	}

	// 计算文件的hash
	fileContent.Seek(0, 0) // 重置文件指针
	hash := util.CalculateFileHash(fileContent)
	if hash == "" {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "计算文件hash失败", Data: "null"})
		return
	}

	// 查询文件是否已存在
	var uploadFile mydb.StructUploadFile
	params := mydb.QueryParams{
		Condition: fmt.Sprintf("hash='%s'", hash),
	}
	existingFile, err := uploadFile.Find(params)
	if err == nil {
		c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "文件已存在", Data: map[string]interface{}{"file_path": existingFile.FilePath}})
		return
	}

	// 获取当前日期
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	// 使用哈希值作为文件名
	fileName := hash + ext

	// 构建存储路径
	filePath := filepath.Join("uploads", "img", year, month, day, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "保存文件失败", Data: "null"})
		return
	}

	// 插入文件记录
	newFile := mydb.StructUploadFile{
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   file.Size,
		Hash:       hash,
		FileType:   "img",
		Extension:  ext,
		UploadTime: util.GetTimestamp(10),
	}
	_, _, err = newFile.Insert([]mydb.StructUploadFile{newFile})
	if err != nil {
		fmt.Println(err)
		util.WrapError(err, "插入文件记录失败")
		c.JSON(http.StatusInternalServerError, util.APIResponse{Code: http.StatusInternalServerError, Message: "插入文件记录失败", Data: "null"})
		return
	}
	img_return_data := ImgReturnData{Hash: hash, ImgPath: "/images/" + fileName}
	c.JSON(http.StatusOK, util.APIResponse{Code: http.StatusOK, Message: "文件上传成功", Data: img_return_data})
}
