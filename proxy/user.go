package proxy

import (
	"io"
	"net/http"
	"time"

	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

func CreateUser(url string, username string, body io.Reader, appG *app.Gin) {
	peer, ok := httpPool.PickPeer(username)
	if !ok {
		logger.GetInstance().WithField("username", username).Errorln("no peers found")
		appG.Response(http.StatusInternalServerError, constval.UnknownError, nil)
		return
	}

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(peer.Addr()+url, "application/json;charset=utf-8", body)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"url": peer.Addr() + url,
			"err": err,
		}).Errorln("request peer error")
		appG.Response(http.StatusBadGateway, constval.UnknownError, nil)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("read resp body error")
		appG.Response(http.StatusInternalServerError, constval.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, constval.OK, map[string]interface{}{"user_id": string(data)})
}
