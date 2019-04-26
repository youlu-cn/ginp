package ginp

import (
	"github.com/gin-gonic/gin"
	"errors"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

var (
	checkToken   = true
	tokenSignKey = "default_key"
	Engine       = gin.Default()
)

func Mode() string {
	return gin.Mode()
}

func SetMode(mode string) {
	gin.SetMode(mode)
}

func SetEngine(engine *gin.Engine) {
	Engine = engine
}

func SetTokenSignKey(signKey string) {
	tokenSignKey = signKey
}

func DisableToken() {
	checkToken = false
}

func AddController(ctl Controller) error {
	if ctl.Group() == "" {
		return errors.New("invalid group name")
	}
	rg := Engine.Group("/" + ctl.Group())
	rg.Any("/*action", Auth(ctl), handleRequest(ctl))
	return nil
}
