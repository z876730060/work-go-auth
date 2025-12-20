package user

import (
	"errors"
	"fmt"
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

func NewHandler(baseHandler common.BaseHandler) *Handler {
	return &Handler{
		l:           baseHandler.Log.With(common.HANDLER, "userHandler"),
		db:          baseHandler.DB,
		redisClient: baseHandler.RedisClient,
		info:        baseHandler.Info,
	}
}

func (h *Handler) Register(e *gin.Engine) {
	e.POST("/user/list", h.List)
	e.POST("/user", h.Add)
	e.DELETE("/user/:id", h.Del)
	e.GET("/user/:id", h.GetDetail)
	e.PUT("/user", h.Update)
	e.POST("/user/role", h.BindRole)
	e.GET("/user/role/:id", h.GetRole)
	e.GET("/user/me", h.Me)
}

func (h *Handler) List(c *gin.Context) {
	type reqBody struct {
		common.Page
		ID       string `json:"ID"`
		Username string `json:"username"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	var rbody reqBody
	err := c.Bind(&rbody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.l.Info("", "reqBody", rbody)

	var datas []User
	var count int64
	query := h.db.Model(&User{})
	if rbody.ID != "" {
		query = query.Where("id = ?", rbody.ID)
	}
	if rbody.Username != "" {
		query = query.Where("username LIKE ?", "%"+rbody.Username+"%")
	}
	if rbody.Fullname != "" {
		query = query.Where("fullname LIKE ?", "%"+rbody.Fullname+"%")
	}
	if rbody.Email != "" {
		query = query.Where("email LIKE ?", "%"+rbody.Email+"%")
	}
	if rbody.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+rbody.Phone+"%")
	}
	err = query.Count(&count).Order("id").Limit(rbody.Size).Offset((rbody.Page.Page - 1) * rbody.Size).Find(&datas).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"records": datas,
			"total":   count,
		},
		"info": h.info,
	})
}

func (h *Handler) Add(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	//校验用户名是否存在
	var count int64
	h.db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, common.RespErr("username already exists", h.info))
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common.RespErr("username already exists", h.info))
			return
		} else {
			c.JSON(http.StatusInternalServerError, common.RespErr("create user failed", h.info))
			return
		}
	}

	h.l.Info("Add user", "user", user)

	c.JSON(http.StatusOK, common.RespOk("create user success", nil, h.info))
}

func (h *Handler) Del(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid id", h.info))
		return
	}

	if err := h.db.Delete(&User{Model: gorm.Model{ID: uid}}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr("delete user failed", h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("delete user success", nil, h.info))
}

func (h *Handler) Update(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	h.l.Info("Update user", "user", user)

	// 校验用户名是否存在
	var count int64
	h.db.Model(&User{}).Where("username = ? AND id != ?", user.Username, user.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, common.RespErr("username already exists", h.info))
		return
	}

	// 校验密码是否为空
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, common.RespErr("password cannot be empty", h.info))
		return
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr("update user failed", h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("update user success", nil, h.info))
}

func (h *Handler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid id", h.info))
		return
	}

	var user User
	if err := h.db.Where(User{Model: gorm.Model{ID: uid}}).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr("get user detail failed", h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get user detail success", user, h.info))
}

func (h *Handler) BindRole(c *gin.Context) {
	var reqBody struct {
		ID       string   `json:"ID"`
		RoleKeys []string `json:"roleKeys"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	uid, err := common.ParseID(reqBody.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid id", h.info))
		return
	}

	h.l.Info("BindRole reqBody", "reqBody", reqBody)
	tx := h.db.Begin()
	defer tx.Rollback()

	// 删除用户角色
	if err := tx.Where(UserRole{UserID: uid}).Delete(&UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr("bind role failed", h.info))
		return
	}

	// 绑定用户角色
	for _, roleKey := range reqBody.RoleKeys {
		roleID, err := common.ParseID(roleKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.RespErr("invalid role key", h.info))
			return
		}
		if err := tx.Create(&UserRole{UserID: uid, RoleID: roleID}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, common.RespErr("bind role failed", h.info))
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, common.RespOk("bind role success", nil, h.info))
}

func (h *Handler) GetRole(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid id", h.info))
		return
	}

	var roles []UserRole
	if err := h.db.Where(UserRole{UserID: uid}).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr("get role failed", h.info))
		return
	}

	// 转换为角色键值对
	roleKeys := make([]string, 0)
	for _, role := range roles {
		roleKeys = append(roleKeys, fmt.Sprintf("%d", role.RoleID))
	}

	c.JSON(http.StatusOK, common.RespOk("get role success", roleKeys, h.info))
}

func (h *Handler) Me(c *gin.Context) {

}
