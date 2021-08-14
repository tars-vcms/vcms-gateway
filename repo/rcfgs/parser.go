package rcfgs

import (
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars/util/conf"
	"reflect"
)

var parserMap = map[CfgType]ParseFunc{
	JSON:   parseJSON,
	YAML:   parseYAML,
	TEXT:   parseText,
	STRUCT: parseStruct,
}

func getParseFunc(cfgType CfgType) ParseFunc {
	return parserMap[cfgType]
}

func parseStruct(from string, s interface{}) error {
	cfg := conf.New()
	if err := cfg.InitFromString(from); err != nil {
		return err
	}
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		path := tf.Tag.Get("tars")
		switch tf.Type.Kind() {
		case reflect.String:
			vf.SetString(cfg.GetString(path))
			break
		case reflect.Map:
			for k, v := range cfg.GetMap(path) {
				vf.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
			break
		}
	}
	return nil
}

func parseJSON(from string, dest interface{}) error {
	return json.Unmarshal([]byte(from), dest)
}

func parseYAML(from string, dest interface{}) error {
	return json.Unmarshal([]byte(from), dest)
}

func parseText(from string, dest interface{}) error {
	*(dest).(*string) = from
	return nil
}
