package repository

import "strconv"

func GetString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func GetStringFromIntPointer(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

func GetStringFromInt(i int) string {
	return strconv.Itoa(i)
}

func GetIntFromStringPointer(s *string) int {
	if s == nil {
		return 0
	}
	i, _ := strconv.Atoi(*s)
	return i
}

func GetIntFromString(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
