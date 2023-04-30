package utils

import (
	"encoding/json"
	"github.com/jinzhu/copier"
)

func DeepCopyByJson(src interface{}, dst interface{}) error {
	if tmp, err := json.Marshal(src); err != nil {
		return err
	} else {
		err := json.Unmarshal(tmp, dst)
		return err
	}
}
func DeepCopy(from interface{}, to interface{}) error {
	if err := copier.CopyWithOption(to, from, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}
	return nil
}
