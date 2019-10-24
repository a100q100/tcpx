package configService

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/fwhezfwhez/model_convert/util"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"tcpx/projects/jelly/control/config/configModel"
	"tcpx/projects/jelly/dependency/db"
)

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateListAPI().
func HTTPListConfig(c *gin.Context) {
	var engine = db.DB.Model(&configModel.Config{})
	id := c.DefaultQuery("id", "")
	if id != "" {
		engine = engine.Where("id != ?", id)
	}
	configId := c.DefaultQuery("config_id", "")
	if configId != "" {
		engine = engine.Where("config_id != ?", configId)
	}
	env := c.DefaultQuery("env", "")
	if env != "" {
		engine = engine.Where("env != ?", env)
	}

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "20")
	orderBy := c.DefaultQuery("order_by", "")
	var count int
	if e := engine.Count(&count).Error; e != nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	var list = make([]configModel.Config, 0, 20)
	if count == 0 {
		c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
		return
	}
	limit, offset := util.ToLimitOffset(size, page, count)
	engine = engine.Limit(limit).Offset(offset)
	if orderBy != "" {
		engine = engine.Order(util.GenerateOrderBy(orderBy))
	}
	if e := engine.Find(&list).Error; e != nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateGetOneAPI().
func HTTPGetConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, e := strconv.Atoi(id)
	if e != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
		return
	}
	var count int
	if e := db.DB.Model(&configModel.Config{}).Where("id=?", idInt).Count(&count).Error; e != nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
		return
	}
	var instance configModel.Config
	if e := db.DB.Model(&configModel.Config{}).Where("id=?", id).First(&instance).Error; e != nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": instance})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func HTTPAddConfig (c *gin.Context) {
	var param configModel.Config
	if e := c.Bind(&param); e!=nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	if e:=db.DB.Model(&configModel.Config{}).Create(&param).Error; e!=nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": param})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateDeleteOneAPI().
func HTTPDeleteConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, e := strconv.Atoi(id)
	if e!=nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
		return
	}
	var count int
	if e:=db.DB.Model(&configModel.Config{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	if count ==0 {
		c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
		return
	}
	var instance configModel.Config
	if e:=db.DB.Model(&configModel.Config{}).Where("id=?", id).Delete(&instance).Error; e!=nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateUpdateOneAPI().
func HTTPUpdateConfig(c *gin.Context) {
	id := c.Param("id")
	idInt, e := strconv.Atoi(id)
	if e!=nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
		return
	}
	var count int
	if e:=db.DB.Model(&configModel.Config{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	if count ==0 {
		c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
		return
	}
	var param configModel.Config
	if e:=c.Bind(&param);e!=nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	if !util.IfZero(param.Id) {
		c.JSON(400, gin.H{"message": "field 'Id' can't be modified'"})
		return
	}
	if e:=db.DB.Model(&configModel.Config{}).Where("id=?", id).Updates(param).Error; e!=nil {
		log.Println(e)
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
