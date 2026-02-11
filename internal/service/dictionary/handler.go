package dictionary

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/z876730060/auth/internal/service/common"
	"gorm.io/gorm"
)

type Handler struct {
	l           *slog.Logger
	db          *gorm.DB
	redisClient *redis.Client
	info        common.Info
}

func NewHandler(base common.BaseHandler) *Handler {
	return &Handler{
		l:           base.Log.With(common.HANDLER, "dictionaryHandler"),
		db:          base.DB,
		redisClient: base.RedisClient,
		info:        base.Info,
	}
}

func (h *Handler) Register(e *gin.Engine) {
	e.POST("/dictionary/list", h.List)
	e.POST("/dictionary/subList", h.SubList)
	e.GET("/dictionary/get", h.Get)
	e.POST("/dictionary/add", h.Add)
	e.PUT("/dictionary/update", h.Update)
	e.POST("/dictionary/addSubItem", h.AddSubItem)
	e.PUT("/dictionary/updateSubItem", h.UpdateSubItem)
	e.DELETE("/dictionary/delete", h.Delete)
	e.DELETE("/dictionary/deleteSubItem", h.DeleteSubItem)
}

func (h *Handler) List(c *gin.Context) {
	var req struct {
		Page common.Page
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var dictionaries []Dictionary
	if err := h.db.Preload("dictionary_item").Offset((req.Page.Page - 1) * req.Page.Size).Limit(req.Page.Size).Find(&dictionaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}
	c.JSON(http.StatusOK, common.RespOk("get dictionary list success", dictionaries, h.info))
}

func (h *Handler) SubList(c *gin.Context) {
	var subItems []DictionaryItem
	dictKey := c.Query("key")
	if dictKey == "" {
		c.JSON(http.StatusBadRequest, common.RespErr("key is required", h.info))
		return
	}
	var dictionary Dictionary
	if err := h.db.Where("key = ?", dictKey).First(&dictionary).Error; err != nil {
		c.JSON(http.StatusNotFound, common.RespErr(err.Error(), h.info))
		return
	}
	if err := h.db.Where("dictionary_id = ?", dictionary.ID).Find(&subItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}
	c.JSON(http.StatusOK, common.RespOk("get sub item list success", subItems, h.info))
}

func (h *Handler) Get(c *gin.Context) {
	dictKey := c.Query("key")
	if dictKey == "" {
		c.JSON(http.StatusBadRequest, common.RespErr("key is required", h.info))
		return
	}

	var dictionary Dictionary
	if err := h.db.Where("key = ?", dictKey).Preload("dictionary_item").First(&dictionary).Error; err != nil {
		c.JSON(http.StatusNotFound, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get dictionary success", dictionary, h.info))
}

func (h *Handler) Add(c *gin.Context) {
	var body Dictionary
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// 检查是否存在相同的key
	var count int64
	h.db.Model(&Dictionary{}).Where("key = ?", body.Key).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Add dictionary", "id", body.ID)

	c.JSON(http.StatusOK, common.RespOk("add dictionary success", body, h.info))
	if err := h.db.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}
	h.l.Info("Add dictionary", "id", body.ID)
	c.JSON(http.StatusOK, common.RespOk("add dictionary success", body, h.info))
}

func (h *Handler) Update(c *gin.Context) {
	var body Dictionary
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// 检查是否存在相同的key
	var count int64
	h.db.Model(&Dictionary{}).Where("key = ? AND id <> ?", body.Key, body.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Save(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Update dictionary", "id", body.ID)

	c.JSON(http.StatusOK, common.RespOk("update dictionary success", body, h.info))
}

func (h *Handler) AddSubItem(c *gin.Context) {
	var body DictionaryItem
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}
	// 检查是否存在相同的key
	var count int64
	h.db.Model(&DictionaryItem{}).Where("key = ? AND dictionary_id = ?", body.Key, body.DictionaryID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Add sub item", "id", body.ID)

	c.JSON(http.StatusOK, common.RespOk("add sub item success", body, h.info))
}

func (h *Handler) UpdateSubItem(c *gin.Context) {
	var body DictionaryItem
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}
	// 检查是否存在相同的key
	var count int64
	h.db.Model(&DictionaryItem{}).Where("key = ? AND dictionary_id = ? AND id <> ?", body.Key, body.DictionaryID, body.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Save(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Update sub item", "id", body.ID)

	c.JSON(http.StatusOK, common.RespOk("update sub item success", body, h.info))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var dictionary Dictionary
	if err := h.db.Where("id = ?", uid).First(&dictionary).Error; err != nil {
		c.JSON(http.StatusNotFound, common.RespErr(err.Error(), h.info))
		return
	}

	if err := h.db.Delete(&dictionary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Delete dictionary", "id", uid)

	c.JSON(http.StatusOK, common.RespOk("delete dictionary success", nil, h.info))
}

func (h *Handler) DeleteSubItem(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var subItem DictionaryItem
	if err := h.db.Where("id = ?", uid).First(&subItem).Error; err != nil {
		c.JSON(http.StatusNotFound, common.RespErr(err.Error(), h.info))
		return
	}

	if err := h.db.Delete(&subItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Delete sub item", "id", uid)

	c.JSON(http.StatusOK, common.RespOk("delete sub item success", nil, h.info))
}
