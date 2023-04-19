package initialize

import (
	"chatgpt-web-go/src/global"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

func ViperInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./src")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	err := viper.Unmarshal(&global.Cfg, viper.DecodeHook(intToTimeDurationSecondHookFunc))
	if err != nil {
		panic(err)
	}
}

func intToTimeDurationSecondHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.Int {
		return data, nil
	}
	if t == reflect.TypeOf(time.Duration(5)) {
		return time.Duration(data.(int)) * time.Second, nil
	}

	// Convert it by parsing
	return data, nil
}
