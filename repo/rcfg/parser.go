package rcfg

import "encoding/json"

var parserMap = map[CfgType]ParseFunc{
	JSON: parseJSON,
	YAML: parseYAML,
	TEXT: parseText,
}

func getParseFunc(cfgType CfgType) ParseFunc {
	return parserMap[cfgType]
}

func parseJSON(content string, dest interface{}) error {
	return json.Unmarshal([]byte(content), dest)
}

func parseYAML(content string, dest interface{}) error {
	return json.Unmarshal([]byte(content), dest)
}

func parseText(content string, dest interface{}) error {
	*(dest).(*string) = content
	return nil
}
