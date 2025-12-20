package role

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
	"gorm.io/gorm"
)

type Handler struct {
	l    *slog.Logger
	info common.Info
	db   *gorm.DB
}

func (h *Handler) Register(e *gin.Engine) {
	e.POST("/role/list", h.List)
	e.POST("/role", h.Add)
	e.GET("/role/:id", h.GetDetail)
	e.PUT("/role", h.Update)
	e.DELETE("/role/:id", h.Del)
	e.GET("/role/tree", h.GetTree)
}

func NewHandler(baseHandler common.BaseHandler) *Handler {
	return &Handler{
		l:    baseHandler.Log.With(common.HANDLER, "roleHandler"),
		info: baseHandler.Info,
		db:   baseHandler.DB,
	}
}

func (h *Handler) List(c *gin.Context) {
	type rBody struct {
		common.Page
	}
	var req rBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	var data []Role
	var count int64
	h.db.Model(&Role{}).Count(&count).Order("id").Offset((req.Page.Page - 1) * req.Size).Limit(req.Size).Find(&data)

	// 转换为响应格式
	c.JSON(http.StatusOK, common.RespOk("get role list success", gin.H{
		"records": data,
		"total":   count,
	}, h.info))
}

func (h *Handler) Add(c *gin.Context) {
	type rBody struct {
		Name           string   `json:"name"`
		MenuPermission []string `json:"menuPermission"`
	}
	var req rBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	// 首先创建角色
	role := Role{
		Name: req.Name,
	}
	if err := h.db.Create(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common.RespErr("role name already exists", h.info))
			return
		} else {
			c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
			return
		}
	}

	h.l.Info("Add role", "role", role)

	// 然后创建角色菜单关系
	for _, menuKey := range req.MenuPermission {
		RoleMenu := RoleMenu{
			Rid:     role.ID,
			MenuKey: menuKey,
		}
		if err := h.db.Create(&RoleMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, common.RespErr("menu permission already exists", h.info))
				return
			} else {
				c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
				return
			}
		}
	}

	c.JSON(http.StatusOK, common.RespOk("create role success", nil, h.info))
}
func (h *Handler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var role Role
	if err := h.db.Where("id = ?", uid).First(&role).Error; err != nil {
		// 区分记录不存在和其他数据库错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.RespErr("role not found", h.info))
			return
		}
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	var RoleMenu []RoleMenu
	if err := h.db.Where("rid = ?", role.ID).Find(&RoleMenu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	var menuPermission []string
	for _, table := range RoleMenu {
		menuPermission = append(menuPermission, table.MenuKey)
	}

	// 转换为响应格式
	c.JSON(http.StatusOK, common.RespOk("get role detail success", gin.H{
		"role":           role,
		"menuPermission": menuPermission,
	}, h.info))
}

func (h *Handler) Update(c *gin.Context) {
	type rBody struct {
		ID             string   `json:"ID"`
		Name           string   `json:"name"`
		Comment        string   `json:"comment"`
		MenuPermission []string `json:"menuPermission"`
	}
	var req rBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	uid, err := common.ParseID(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// 首先更新角色
	role := Role{
		Model: gorm.Model{
			ID: uid,
		},
		Name:    req.Name,
		Comment: req.Comment,
	}
	if err := h.db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	// 然后删除旧的角色菜单关系
	if err := h.db.Where("rid = ?", uid).Unscoped().Delete(&RoleMenu{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	// 最后创建新的角色菜单关系
	for _, menuKey := range req.MenuPermission {
		RoleMenu := RoleMenu{
			Rid:     uid,
			MenuKey: menuKey,
		}
		if err := h.db.Create(&RoleMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, common.RespErr("menu permission already exists", h.info))
				return
			} else {
				c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
				return
			}
		}
	}

	c.JSON(http.StatusOK, common.RespOk("update role success", nil, h.info))
}
func (h *Handler) Del(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// 首先删除角色菜单关系
	if err := h.db.Where("rid = ?", uid).Delete(&RoleMenu{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	// 然后删除角色
	if err := h.db.Delete(&Role{}, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("delete role success", nil, h.info))
}

func (h *Handler) GetTree(c *gin.Context) {
	var roles []Role
	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	var roleTree []RoleTree
	for _, role := range roles {
		roleTree = append(roleTree, RoleTree{
			Title:    role.Name,
			Key:      fmt.Sprintf("%d", role.ID),
			Children: []*RoleTree{},
		})
	}

	c.JSON(http.StatusOK, common.RespOk("get role tree success", roleTree, h.info))
}
