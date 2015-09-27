package utils

import (
	"encoding/json"
)

func JsonStringify(item interface{}) (string, error) {
	ret, err := json.Marshal(item)
	return string(ret), err
}

func JsonParse(str string, item *interface{}) error {
	err := json.Unmarshal([]byte(str), item)
	return err
}
