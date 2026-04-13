package configx

import (
	"strings"

	"github.com/jiny3/gopkg/filex"
	"github.com/spf13/viper"
)

// set input viper config
func load(v *viper.Viper, path string) error {
	if path != "" {
		args := filex.Parse(path)
		v.SetConfigName(args.Name)
		v.SetConfigType(args.Type)
		v.AddConfigPath(args.Dir)
	}

	// env 覆盖
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v.ReadInConfig()
}
