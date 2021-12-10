package log

//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"debug":       false,
		"processId":   "",
		"service":     "mqant",
		"logDir":      "",
		"logSettings": map[string]interface{}{},
		"biDir":       "",
		"biSettings":  map[string]interface{}{},
		"useBi":       true,
		"useLog":      true,
	}
}
