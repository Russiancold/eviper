package eviper

import (
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type EViper struct {
	*viper.Viper
}

func New(v *viper.Viper) *EViper {
	return &EViper{v}
}

func (e *EViper) Unmarshal(rawVal interface{}) error {
	if err := e.Viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// 	do nothing
		default:
			return err
		}
	}

	_ = e.Viper.Unmarshal(rawVal)
	e.readEnvs(rawVal)
	return e.Viper.Unmarshal(rawVal)
}

func (e *EViper) readEnvs(rawVal interface{}) {
	e.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	e.bindEnvs(rawVal)
}

func (e *EViper) bindEnvs(in interface{}, prev ...string) {
	ifv := reflect.ValueOf(in)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}

	for i := 0; i < ifv.NumField(); i++ {
		fv := ifv.Field(i)
		t := ifv.Type().Field(i)
		tv, ok := t.Tag.Lookup("env")
		if ok {
			if tv == ",squash" {
				e.bindEnvs(fv.Interface(), prev...)
				continue
			}
		}

		env := strings.Join(append(prev, t.Name), ".")
		switch fv.Kind() {
		case reflect.Struct:
			e.bindEnvs(fv.Interface(), append(prev, t.Name)...)
		case reflect.Slice:
			e.Viper.SetTypeByDefaultValue(true)
			e.Viper.SetDefault(env, []string{})
			_ = e.Viper.BindEnv(env, tv)
		default:
			_ = e.Viper.BindEnv(env, tv)
		}
	}
}
