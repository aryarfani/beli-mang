package str

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ShowString ...
func ShowString(isShow bool, data string) string {
	if isShow {
		return data
	}

	return ""
}

// EmptyString ...
func EmptyString(text string) *string {
	if text == "" {
		return nil
	}
	return &text
}

// EmptyInt ...
func EmptyInt(number int) *int {
	if number == 0 {
		return nil
	}
	return &number
}

// StringToInt ...
func StringToInt(data string) int {
	res, err := strconv.Atoi(data)
	if err != nil {
		res = 0
	}

	return res
}

// StringToFloat ...
func StringToFloat(data string, bitize int) float64 {
	res, err := strconv.ParseFloat(data, bitize)
	if err != nil {
		res = 0
	}

	return res
}

// StringToBool ...
func StringToBool(data string) bool {
	res, err := strconv.ParseBool(data)
	if err != nil {
		res = false
	}

	return res
}

// StringToBoolString ...
func StringToBoolString(data string) string {
	_, err := strconv.ParseBool(data)
	if err != nil {
		return "false"
	}

	return data
}

// RandomString ...
func RandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}

// Unique ...
func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// CheckEmail ...
func CheckEmail(text string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(text)
}

// IsValidUUID ...
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// TimeToString ...
func TimeToString(data time.Time, format string) string {
	if data.IsZero() {
		return ""
	}

	return data.Format(format)
}
