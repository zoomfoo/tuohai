package mainsite

import (
	"gopkg.in/gin-gonic/gin.v1"
)

//主站router
func NewMainSiteRouter(rg *gin.RouterGroup) {
	//创建项目群
	rg.POST("/project/groups", CreateProjectGroup())
	rg.POST("/team/groups", CreateTeamGroup())
}
