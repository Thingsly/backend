// serialize.go is a Go file that contains the code to serialize data to JSON format.

package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func SerializeData(source, target interface{}) (interface{},error )  {
	jsonData, err := json.Marshal(source)
	if err != nil {
		logrus.Error("JSON serialization failed:", err)
		return nil,err
	}

	err = json.Unmarshal(jsonData, &target)
	if err != nil {
		logrus.Error("JSON deserialization failed:", err)
		return nil,err
	}

	return target,nil
}