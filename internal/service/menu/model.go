package menu

import (
	"github.com/z876730060/auth/pkg/menu"
	"gorm.io/gorm"
)

type MenuTable struct {
	gorm.Model
	menu.Menu
	menu.Route
	OrderId      int      `json:"orderId" gorm:"default:0"`
	MicroAppBean MicroApp `json:"microAppBean" gorm:"foreignKey:MicroApp;references:Key"`
}

type TreeMenu struct {
	Title    string      `json:"title"`
	Key      string      `json:"key"`
	Children []*TreeMenu `json:"children"`
}

func (m *MenuTable) TableName() string {
	return "menu"
}

type MicroApp struct {
	gorm.Model
	Name    string `json:"name"`
	Key     string `json:"key"`
	BaseUrl string `json:"baseUrl"`
}

func (m *MicroApp) TableName() string {
	return "micro_app"
}

func InitMenuTable(db *gorm.DB) {
	db.AutoMigrate(&MenuTable{})
	db.AutoMigrate(&MicroApp{})

	var count int64
	db.Model(&MenuTable{}).Count(&count)
	if count > 0 {
		return
	}

	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:   "/",
			Label: "首页",
		},
		Route: menu.Route{
			Path:      "/",
			Component: "page/board/Board",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:   "/user",
			Label: "用户管理",
		},
		Route: menu.Route{
			Path:      "/user",
			Component: "page/user/User",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:   "/role",
			Label: "角色管理",
		},
		Route: menu.Route{
			Path:      "/role",
			Component: "page/role/Role",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:   "/menu",
			Label: "菜单管理",
		},
		Route: menu.Route{
			Path:      "/menu",
			Component: "page/menu/Menu",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/menu/add",
			Label:     "添加菜单",
			ParentKey: "/menu",
		},
		Route: menu.Route{
			Path:      "/menu/add",
			Component: "page/menu/AddMenu",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/role/add",
			Label:     "添加角色",
			ParentKey: "/role",
		},
		Route: menu.Route{
			Path:      "/role/add",
			Component: "page/role/AddRole",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/user/add",
			Label:     "添加用户",
			ParentKey: "/user",
		},
		Route: menu.Route{
			Path:      "/user/add",
			Component: "page/user/AddUser",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/menu/edit",
			Label:     "编辑菜单",
			ParentKey: "/menu",
		},
		Route: menu.Route{
			Path:      "/menu/edit",
			Component: "page/menu/EditMenu",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/role/edit",
			Label:     "编辑角色",
			ParentKey: "/role",
		},
		Route: menu.Route{
			Path:      "/role/edit",
			Component: "page/role/EditRole",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/user/edit",
			Label:     "编辑用户",
			ParentKey: "/user",
		},
		Route: menu.Route{
			Path:      "/user/edit",
			Component: "page/user/EditUser",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/user/role",
			Label:     "用户角色",
			ParentKey: "/user",
		},
		Route: menu.Route{
			Path:      "/user/role",
			Component: "page/user/UserRole",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:   "/work-vue",
			Label: "vue 应用",
		},
		Route: menu.Route{
			Path:      "/work-vue",
			Component: "page/work-vue/WorkVue",
			Other:     true,
			MicroApp:  "work-vue",
		},
	})
	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/work-vue/monitor",
			Label:     "vue 应用 - 监控看板",
			ParentKey: "/work-vue",
		},
		Route: menu.Route{
			Path:      "/work-vue/monitor",
			Component: "page/work-vue/WorkVue",
			Other:     true,
			MicroApp:  "work-vue",
		},
	})

	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/work-vue/board",
			Label:     "vue 应用 - 首页",
			ParentKey: "/work-vue",
		},
		Route: menu.Route{
			Path:      "/work-vue/board",
			Component: "page/work-vue/WorkVue",
			Other:     true,
			MicroApp:  "work-vue",
		},
	})

	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/menu/micro-app",
			Label:     "微应用",
			ParentKey: "/menu",
		},
		Route: menu.Route{
			Path:      "/menu/micro-app",
			Component: "page/menu/microApp/MicroApp",
			Other:     false,
			MicroApp:  "",
		},
	})

	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/menu/micro-app/add",
			Label:     "添加微应用",
			ParentKey: "/menu/micro-app",
		},
		Route: menu.Route{
			Path:      "/menu/micro-app/add",
			Component: "page/menu/microApp/AddMicroApp",
			Other:     false,
			MicroApp:  "",
		},
	})

	db.Create(&MenuTable{
		Menu: menu.Menu{
			Key:       "/menu/micro-app/edit",
			Label:     "编辑微应用",
			ParentKey: "/menu/micro-app",
		},
		Route: menu.Route{
			Path:      "/menu/micro-app/edit",
			Component: "page/menu/microApp/EditMicroApp",
			Other:     false,
			MicroApp:  "",
		},
	})

	db.Create(&MicroApp{
		Name:    "work-vue",
		Key:     "work-vue",
		BaseUrl: "http://localhost:5174",
	})
}
