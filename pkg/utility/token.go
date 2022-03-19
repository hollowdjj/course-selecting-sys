package utility

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

//time + md5生成随机token
func GenerateToken() string {
	currTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(currTime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
