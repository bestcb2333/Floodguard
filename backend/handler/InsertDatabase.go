package handler

import (
	"fmt"

	"github.com/bestcb2333/FloodGuard/database"
	"github.com/bestcb2333/FloodGuard/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var insertActionMap = map[string]any{
	"add-notice":     database.Notice{},
	"add-floodevent": database.FloodEvent{},
	"add-region":     database.Region{},
}

func InsertDatabase(c *gin.Context) {

	user, err := VerifyJwt(c)
	if err != nil {
		util.Error(c, 400, "用户身份验证失败", nil)
		return
	}

	path := c.Request.URL.Path[1:]
	action := insertActionMap[path]

	if !viper.GetBool(fmt.Sprintf("user-permission.%s", path)) {
		if !user.Admin {
			util.Error(c, 400, "你不是管理员，无法执行此操作", nil)
			return
		}
	}

	if err := c.BindJSON(&action); err != nil {
		util.Error(c, 400, "无法读取你的请求", err)
		return
	}

	if err := DB.Create(&action).Error; err != nil {
		util.Error(c, 500, "数据上传失败", err)
		return
	}

	util.Info(c, 200, "操作成功", nil)
}
