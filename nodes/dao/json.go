package dao

import (
	"fmt"
	"strconv"
)

type JSON map[string]interface{}

func (json JSON) GetString(name string) (string, bool) {
	value, has := json[name]
	if has {
		return fmt.Sprintf("%v", value), has
	} else {
		return "", has
	}
}

func (json JSON) String(name string) string {
	value, _ := json.GetString(name)
	return value
}

func (json JSON) GetInt(name string) (int, bool) {
	value, has := json[name]
	if !has {
		return 0, false
	} else if intValue, match := value.(int); match {
		return intValue, true
	} else {
		if i, err := strconv.Atoi(fmt.Sprintf("%v", value)); err != nil {
			return 0, false
		} else {
			return i, true
		}
	}
}

func (json JSON) Int(name string, def int) int {
	value, has := json.GetInt(name)
	if !has {
		return def
	}
	return value
}
