package utils

import "encoding/json"

func DeepCopyByJson(src interface{}, dst interface{}) error {
	if tmp, err := json.Marshal(src); err != nil {
		return err
	} else {
		err := json.Unmarshal(tmp, dst)
		return err
	}
}
