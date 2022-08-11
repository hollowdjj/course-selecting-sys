package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"github.com/sirupsen/logrus"
)

var ServerMode = "usual"

type SwitchModeForm struct {
	Mode string `valid:"Required"`
}

//used for switching mode, only effective on main host
func SwitchMode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form SwitchModeForm
	)
	currHost := conf.GetApp().Host
	mainHost := conf.GetApp().MainHost
	if currHost != mainHost {
		logger.GetInstance().WithFields(logrus.Fields{
			"curr host": currHost,
			"main host": mainHost,
		}).Errorln("only main host can switch host")
		appG.Response(http.StatusBadRequest, constval.PermDenied, nil)
		return
	}

	//form validation
	funcs := app.CustomFunc{
		"ModeCheck": utility.ModeCheck,
	}
	httpCode, errCode := app.BindAndValidCustom(c, &form, funcs, true)
	if errCode != constval.OK {
		logger.GetInstance().WithField("mode", form.Mode).Errorln("switch mode form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//switch mode
	ServerMode = form.Mode
}
