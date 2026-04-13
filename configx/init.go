package configx

import (
	"github.com/sirupsen/logrus"
)

func init() {
	defaultPath := "config/global.toml"
	err := Load(defaultPath)
	if err != nil {
		logrus.WithField("path", defaultPath).Debug("default config not found")
	}
}
