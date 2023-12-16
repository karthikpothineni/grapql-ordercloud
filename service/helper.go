package service

func GetString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func GetBool(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}

func GetBoolPointer(val bool) *bool {
	return &val
}
