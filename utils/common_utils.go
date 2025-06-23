package utils

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnUsed(x ...interface{}) {}

func RemoveDuplicateObjectID(slice []primitive.ObjectID) []primitive.ObjectID {
	if len(slice) < 2 {
		return slice
	}

	allKeys := make(map[primitive.ObjectID]bool)
	list := []primitive.ObjectID{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateString(slice []string) []string {
	if len(slice) < 2 {
		return slice
	}

	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateAndEmptyString(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if item != "" {
			if _, value := allKeys[item]; !value {
				allKeys[item] = true
				list = append(list, item)
			}
		}
	}
	return list
}

func RemoveFromObjectIdSlice(slice []primitive.ObjectID, value primitive.ObjectID) []primitive.ObjectID {
	result := []primitive.ObjectID{}
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

func RemoveFromAnySlice(slice []any, value any) []any {
	result := []any{}
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

func HasDuplicateObjectIds(slice []primitive.ObjectID) bool {
	if len(slice) < 2 {
		return false
	}

	return !(len(slice) == len(RemoveDuplicateObjectID(slice)))
}

func HasDuplicateStrings(slice []string) bool {
	if len(slice) < 2 {
		return false
	}

	return !(len(slice) == len(RemoveDuplicateString(slice)))
}

func ConvertInterfaceToStruct(data any, object interface{}) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid data format")
	}

	// Marshal the map back to JSON
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &object)
	if err != nil {
		return err
	}

	if err := validate.Struct(object); err != nil {
		return err
	}

	return nil
}
