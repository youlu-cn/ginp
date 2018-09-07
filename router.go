package ginp

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

var (
	Engine = gin.Default()
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

func AddController(ctl Controller) error {
	if ctl.Group() == "" {
		return errors.New("invalid group name")
	}
	rg := Engine.Group("/" + ctl.Group())
	rg.Any("/*action", handleRequest(ctl))
	return nil
}
