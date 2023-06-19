package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiKeyHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewApiKeyHandler(app *core.AppServer, db *gorm.DB) *ApiKeyHandler {
	h := ApiKeyHandler{db: db}
	h.App = app
	return &h
}

func (h *ApiKeyHandler) Add(c *gin.Context) {
	var data struct {
		Key string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 获取当前登录用户
	var userId uint = 0
	user, err := utils.GetLoginUser(c, h.db)
	if err == nil {
		userId = user.Id
	}
	var key = model.ApiKey{Value: data.Key, UserId: userId}
	res := h.db.Create(&key)
	if res.Error != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "操作失败"})
		return
	}
	resp.SUCCESS(c, key)
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	offset := (page - 1) * pageSize
	var items []model.ApiKey
	var keys = make([]vo.ApiKey, 0)
	var total int64
	h.db.Model(&model.ApiKey{}).Count(&total)
	res := h.db.Offset(offset).Limit(pageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var key vo.ApiKey
			err := utils.CopyObject(item, &key)
			if err == nil {
				key.Id = item.Id
				key.CreatedAt = item.CreatedAt.Unix()
				key.UpdatedAt = item.UpdatedAt.Unix()
				keys = append(keys, key)
			} else {
				logger.Error(err)
			}
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, keys)
	resp.SUCCESS(c, pageVo)
}