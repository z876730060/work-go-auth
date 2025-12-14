package menu

import (
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
	"github.com/z876730060/auth/internal/service/role"
	"github.com/z876730060/auth/pkg/menu"

	"gorm.io/gorm"
)

type Handler struct {
	l    *slog.Logger
	db   *gorm.DB
	info common.Info
}

func NewHandler(l *slog.Logger, db *gorm.DB, info common.Info) *Handler {
	return &Handler{l: l, db: db, info: info}
}

func (h *Handler) Register(e *gin.Engine) {
	e.GET("/menu", h.GetMenu)
	e.GET("/route", h.GetRoute)
	e.POST("/menu/list", h.List)
	e.POST("/menu", h.Add)
	e.DELETE("/menu/:id", h.Del)
	e.GET("/breadcrumb", h.GetBreadcrumb)
	e.GET("/menu/:id", h.GetDetail)
	e.PUT("/menu", h.Update)
	e.GET("/menu/tree", h.GetTree)
}

func (h *Handler) GetMenu(c *gin.Context) {
	roleIDs := c.GetUintSlice("role")
	if len(roleIDs) == 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("role is empty", h.info))
		return
	}
	h.l.Info("", "roleIDs", roleIDs)

	appId := c.GetHeader("MicroAppId")

	datas := make([]menu.Menu, 0)
	h.l.Info("db is running\n")
	var data []MenuTable
	if slices.Contains(roleIDs, 1) {
		h.l.Info("admin role")
		query := h.db
		if appId != "" {
			query = query.Where("parent_key = (?) and other = true and (hidden = false or hidden is null)",
				h.db.Model(&MenuTable{}).Where("micro_app = ? and parent_key = ?", appId, "").Pluck("key", nil),
			)
		} else {
			query = query.Where("parent_key = ? and (hidden = false or hidden is null)", "")
		}
		query.Order("order_id, ID").Find(&data)
	} else {
		query := h.db
		if appId != "" {
			query = query.Where("parent_key = (?) and key IN (?) and other = true and (hidden = false or hidden is null)",
				h.db.Model(&MenuTable{}).Where("micro_app = ? and parent_key = ?", appId, "").Pluck("key", nil),
				h.db.Model(&role.RoleMenu{}).Where("rid IN ?", roleIDs).Pluck("menu_key", nil),
			)
		} else {
			query = query.Where("parent_key = (?) and key IN (?) and (hidden = false or hidden is null)", "",
				h.db.Model(&role.RoleMenu{}).Where("rid IN ?", roleIDs).Pluck("menu_key", nil),
			)
		}
		query.Order("order_id, ID").Find(&data)
	}

	for _, item := range data {
		datas = append(datas, item.Menu)
	}

	c.JSON(http.StatusOK, common.RespOk("get menu success", datas, h.info))
}

func (h *Handler) GetRoute(c *gin.Context) {
	roleIDs := c.GetUintSlice("role")
	if len(roleIDs) == 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("role is empty", h.info))
		return
	}
	h.l.Info("", "roleIDs: ", roleIDs)

	microAppId := c.GetHeader("MicroAppId")

	var data []MenuTable

	if slices.Contains(roleIDs, 1) {
		h.l.Info("admin role")
		h.db.Find(&data)
	} else {
		query := h.db
		if microAppId != "" {
			query = query.Where("key IN (?) and other = true",
				h.db.Model(&role.RoleMenu{}).Where("rid IN ?", roleIDs).Pluck("menu_key", nil),
			)
		} else {
			query = query.Where("key IN (?)",
				h.db.Model(&role.RoleMenu{}).Where("rid IN ?", roleIDs).Pluck("menu_key", nil),
			)
		}
		query.Find(&data)
	}

	datas := make([]menu.Route, len(data))
	for i, item := range data {
		datas[i] = item.Route
	}

	c.JSON(http.StatusOK, common.RespOk("get route success", datas, h.info))
}

func (h *Handler) List(c *gin.Context) {
	type rbody struct {
		common.Page
		ID        string `json:"ID"`
		Component string `json:"component"`
		Path      string `json:"path"`
		Key       string `json:"key"`
		Label     string `json:"label"`
	}

	var body rbody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("List menu: %v", body)

	var data []MenuTable
	query := h.db.Model(&MenuTable{})
	if body.ID != "" {
		query = query.Where("id = ?", body.ID)
	}
	if body.Component != "" {
		query = query.Where("component LIKE ?", "%"+body.Component+"%")
	}
	if body.Path != "" {
		query = query.Where("path LIKE ?", "%"+body.Path+"%")
	}
	if body.Key != "" {
		query = query.Where("key LIKE ?", "%"+body.Key+"%")
	}
	if body.Label != "" {
		query = query.Where("label LIKE ?", "%"+body.Label+"%")
	}

	var count int64
	if err := query.Preload("MicroAppBean").Count(&count).Order("id").Limit(body.Size).Offset((body.Page.Page - 1) * body.Size).Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("list menu success", gin.H{
		"records": data,
		"total":   count,
	}, h.info))
}

func (h *Handler) Add(c *gin.Context) {
	type rbody struct {
		menu.Menu
		menu.Route
	}

	var body rbody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	menuTable := MenuTable{
		Menu:  body.Menu,
		Route: body.Route,
	}

	// 检查是否存在相同的key
	var count int64
	h.db.Model(&MenuTable{}).Where("key = ?", body.Key).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	// 检查是否存在相同的path
	h.db.Model(&MenuTable{}).Where("path = ?", body.Path).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("path already exists", h.info))
		return
	}

	if err := h.db.Create(&menuTable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Add menu id: %v", menuTable.ID)

	c.JSON(http.StatusOK, common.RespOk("add menu success", menuTable.Menu, h.info))
}

func (h *Handler) Del(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	if err := h.db.Delete(&MenuTable{}, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	h.l.Info("Del menu id: %v", uid)

	c.JSON(http.StatusOK, common.RespOk("del menu success", nil, h.info))
}

func (h *Handler) GetBreadcrumb(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, common.RespErr("path is empty", h.info))
		return
	}

	h.l.Info("Get breadcrumb", "path", path)
	u, err := url.Parse(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}
	path = u.Path

	if u.Query().Has("my-app") {

		appID, err := url.QueryUnescape(u.Query().Get("my-app"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
			return
		}

		path = strings.TrimSuffix(appID, "/")
	}

	datas := make([]string, 0)

	var data MenuTable
	if err := h.db.Where(MenuTable{Menu: menu.Menu{Key: path}}).First(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	datas = getBreadcrumb(h.db, data.Key, datas)

	// 反转数组
	for i, j := 0, len(datas)-1; i < j; i, j = i+1, j-1 {
		datas[i], datas[j] = datas[j], datas[i]
	}

	c.JSON(http.StatusOK, common.RespOk("get breadcrumb success", datas, h.info))
}

func getBreadcrumb(db *gorm.DB, key string, datas []string) []string {
	var data MenuTable
	db.Where(MenuTable{Menu: menu.Menu{Key: key}}).Find(&data)
	datas = append(datas, data.Label)
	if data.ParentKey != "" {
		return getBreadcrumb(db, data.ParentKey, datas)
	}
	return datas
}

func (h *Handler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var data MenuTable
	if err := h.db.Where(MenuTable{Model: gorm.Model{ID: uid}}).First(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get menu detail success", data, h.info))
}

func (h *Handler) Update(c *gin.Context) {
	var menu MenuTable
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	if err := h.db.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("update menu success", menu, h.info))
}

func (h *Handler) GetTree(c *gin.Context) {
	treeData := make([]*TreeMenu, 0)
	var datas []MenuTable
	h.db.Model(&MenuTable{}).Where("parent_key = ''").Find(&datas)

	for _, data := range datas {
		treeData = append(treeData, &TreeMenu{
			Title:    data.Label,
			Key:      data.Key,
			Children: make([]*TreeMenu, 0),
		})
	}

	getTreeMenu(h.db, treeData)

	c.JSON(http.StatusOK, common.RespOk("get menu tree success", treeData, h.info))
}

func getTreeMenu(db *gorm.DB, data []*TreeMenu) {
	for _, item := range data {
		var datas []MenuTable
		db.Model(&MenuTable{}).Where("parent_key = ?", item.Key).Find(&datas)

		for _, data := range datas {
			item.Children = append(item.Children, &TreeMenu{
				Title:    data.Label,
				Key:      data.Key,
				Children: make([]*TreeMenu, 0),
			})
		}
		getTreeMenu(db, item.Children)
	}
}

type MicroAppHandler struct {
	db   *gorm.DB
	l    *slog.Logger
	info common.Info
}

func (h *MicroAppHandler) Register(e *gin.Engine) {
	e.POST("/micro-app/list", h.List)
	e.POST("/micro-app", h.Add)
	e.DELETE("/micro-app/:id", h.Del)
	e.GET("/micro-app/:id", h.GetDetail)
	e.PUT("/micro-app", h.Update)
	e.GET("/micro-app/select", h.GetSelect)
	e.GET("/micro-app/key/:key", h.GetDetailByKey)
}

func NewMicroAppHandler(l *slog.Logger, info common.Info, db *gorm.DB) *MicroAppHandler {
	return &MicroAppHandler{
		db:   db,
		l:    l,
		info: info,
	}
}

func (h *MicroAppHandler) List(c *gin.Context) {
	type rbody struct {
		common.Page
		ID      string `json:"ID"`
		Name    string `json:"name"`
		Key     string `json:"key"`
		BaseUrl string `json:"baseUrl"`
	}

	var body rbody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	offset := (body.Page.Page - 1) * body.Page.Size

	// 构建查询条件
	query := h.db.Model(&MicroApp{})
	if body.ID != "" {
		query = query.Where("ID = ?", body.ID)
	}
	if body.Name != "" {
		query = query.Where("name LIKE ?", "%"+body.Name+"%")
	}
	if body.Key != "" {
		query = query.Where("key LIKE ?", "%"+body.Key+"%")
	}
	if body.BaseUrl != "" {
		query = query.Where("base_url LIKE ?", "%"+body.BaseUrl+"%")
	}

	var datas []MicroApp
	var count int64
	if err := query.Count(&count).Offset(offset).Limit(body.Size).Order("ID").Find(&datas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get micro app list success", gin.H{
		"records": datas,
		"total":   count,
	}, h.info))
}

func (h *MicroAppHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	var data MicroApp
	if err := h.db.Where(MicroApp{Model: gorm.Model{ID: uid}}).First(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get micro app detail success", data, h.info))
}

func (h *MicroAppHandler) Update(c *gin.Context) {
	var app MicroApp
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// key 不能重复
	var count int64
	h.db.Model(&MicroApp{}).Where("key = ? AND id <> ?", app.Key, app.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("update micro app success", app, h.info))
}

func (h *MicroAppHandler) Del(c *gin.Context) {
	id := c.Param("id")

	uid, err := common.ParseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	if err := h.db.Delete(&MicroApp{}, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("delete micro app success", nil, h.info))
}

func (h *MicroAppHandler) Add(c *gin.Context) {
	var app MicroApp
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// key 不能重复
	var count int64
	h.db.Model(&MicroApp{}).Where("key = ?", app.Key).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, common.RespErr("key already exists", h.info))
		return
	}

	if err := h.db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("create micro app success", app, h.info))
}

func (h *MicroAppHandler) GetSelect(c *gin.Context) {
	var datas []MicroApp
	if err := h.db.Order("ID").Find(&datas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	type Select struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	var selects []Select
	for _, data := range datas {
		selects = append(selects, Select{
			Label: data.Name,
			Value: data.Key,
		})
	}

	c.JSON(http.StatusOK, common.RespOk("get micro app select success", selects, h.info))
}

func (h *MicroAppHandler) GetDetailByKey(c *gin.Context) {
	key := c.Param("key")

	var data MicroApp
	if err := h.db.Where(MicroApp{Key: key}).First(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("get micro app detail by key success", data, h.info))
}
