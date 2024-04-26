package api

import (
	"naive-admin-go/db"
	"naive-admin-go/inout"
	"naive-admin-go/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Permissions = &permissions{}

type permissions struct {
}

func (permissions) List(c *gin.Context) {
	var onePermissList = make([]model.Permission, 0)
	db.Dao.Model(model.Permission{}).Where("parentId is NULL").Order("`order` Asc").Find(&onePermissList)
	for i, perm := range onePermissList {
		var twoPerissList []model.Permission
		db.Dao.Model(model.Permission{}).Where("parentId = ?", perm.ID).Where("type = ?", "MENU").Order("`order` Asc").Find(&twoPerissList)
		for i2, perm2 := range twoPerissList {
			var twoPerissList2 []model.Permission
			db.Dao.Model(model.Permission{}).Where("parentId = ?", perm2.ID).Where("type = ?", "MENU").Order("`order` Asc").Find(&twoPerissList2)
			twoPerissList[i2].Children = twoPerissList2
		}
		onePermissList[i].Children = twoPerissList
	}

	Resp.Succ(c, onePermissList)
}

func (permissions) ListButton(c *gin.Context) {
	id := c.Param("id")
	permissions := make([]model.Permission, 0)
	db.Dao.Model(model.Permission{}).Where("parentId = ?", id).Order("`order` Asc").Find(&permissions)
	Resp.Succ(c, permissions)
}

func (permissions) ListPage(c *gin.Context) {
	var data = &inout.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	orm := db.Dao.Model(model.Role{})
	if name != "" {
		orm = orm.Where("name like ?", "%"+name+"%")
	}
	orm.Count(&data.Total)

	orm.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&data.PageData)
	for i, datum := range data.PageData {
		var perIdList []int64
		db.Dao.Model(model.RolePermissionsPermission{}).Where("roleId=?", datum.ID).Select("permissionId").Find(&perIdList)
		data.PageData[i].PermissionIds = perIdList
	}
	Resp.Succ(c, data)
}
func (permissions) Add(c *gin.Context) {
	var params inout.AddPermissionReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	err = db.Dao.Model(model.Permission{}).Create(&model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId, // insert value null
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: IsTrue(params.KeepAlive),
		Show:      params.Show,
		Enable:    params.Enable,
		Order:     params.Order,
	}).Error
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
func (permissions) Delete(c *gin.Context) {
	id := c.Param("id")
	err := db.Dao.Transaction(func(tx *gorm.DB) error {
		tx.Where("id =?", id).Delete(&model.Permission{})
		tx.Where("permissionId =?", id).Delete(&model.RolePermissionsPermission{})
		return nil
	})
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
func (permissions) PatchPermission(c *gin.Context) {
	var params inout.PatchPermissionReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	err = db.Dao.Model(model.Permission{}).Where("id=?", params.Id).Updates(model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId,
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: params.KeepAlive,
		Method:    params.Component,
		Show:      params.Show,
		Enable:    params.Enable,
		Order:     params.Order,
	}).Error
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")

}
func IsTrue(v bool) int {
	if v {
		return 1
	}
	return 0
}
