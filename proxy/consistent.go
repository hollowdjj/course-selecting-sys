package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	httpPool *cache.HttpPool
)

func InitHttpPool() {
	httpPool = cache.NewHttpPool(conf.GetApp().Host)
}

func GetHttpPool() *cache.HttpPool {
	return httpPool
}

//used for registerting consistent hash node
type RegisterHashNodeForm struct {
	Host string `valid:"Required;IP"`
}

//register consistent hash node
func RegisterConsistentHashNode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = &RegisterHashNodeForm{}
	)
	//form validation
	httpCode, errCode := app.BindAndValid(c, form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"slave server": c.ClientIP(),
			"host":         form.Host,
			"msg":          constval.GetErrCodeMsg(errCode),
		}).Infoln("register consistent hash node form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}
	//check whether this server is main server
	if conf.GetApp().Host != conf.GetApp().MainHost {
		logger.GetInstance().WithFields(logrus.Fields{
			"slave server": c.ClientIP(),
			"host":         form.Host,
		}).Errorln("current server is not main server")
		appG.Response(http.StatusBadRequest, constval.PermDenied,
			"must register consistent hash node to main server")
		return
	}
	//register consistent hash node
	httpPool.AddPeers(form.Host)
	logger.GetInstance().WithFields(logrus.Fields{
		"main_host": conf.GetApp().MainHost,
		"curr_host": conf.GetApp().Host,
		"peers":     httpPool.GetPeers(),
	}).Infoln("add peers succ")
	appG.Response(http.StatusOK, constval.OK, nil)
}

//unregister consistent hash node
func UnRegisterConsistentHashNode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = &RegisterHashNodeForm{}
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"slave": c.ClientIP(),
			"host":  form.Host,
			"msg":   constval.GetErrCodeMsg(errCode),
		}).Infoln("unregister consistent hash node form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//check whether this server is main server
	if conf.GetApp().Host != conf.GetApp().MainHost {
		logger.GetInstance().WithFields(logrus.Fields{
			"slave server": c.ClientIP(),
			"host":         form.Host,
		}).Errorln("current server is not main server")
		appG.Response(http.StatusBadRequest, constval.PermDenied,
			"must unregister consistent hash node to main server")
		return
	}

	//unregister
	remainHosts := httpPool.DelPeer(form.Host)
	logger.GetInstance().WithFields(logrus.Fields{
		"slave":        c.ClientIP(),
		"remain_hosts": remainHosts,
	}).Infoln("unregister consistent hash node succ")
	appG.Response(http.StatusOK, constval.OK, map[string]interface{}{"remain_hosts": remainHosts})
}
