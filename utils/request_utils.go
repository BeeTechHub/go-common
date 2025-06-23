package utils

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPaginationParams(c *fiber.Ctx) (int64, int64, int64, int64) {
	pageStr := c.Query("page", "1")
	sizeStr := c.Query("size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size < 1 {
		size = 10
	}
	firstIdx := (page - 1) * size
	lastIdx := firstIdx + size - 1

	return page, size, firstIdx, lastIdx
}

func GetSortParamsForMongoDb(c *fiber.Ctx) bson.D {
	sortStr := c.Query("sort", "")

	sortParams := bson.D{}
	for _, v := range strings.Split(sortStr, ";") {
		sortElement := strings.Split(v, ",")

		if len(sortElement) >= 2 && sortElement[0] != "" {
			if sortElement[1] == "asc" {
				sortParams = append(sortParams, bson.E{sortElement[0], 1})
			} else if sortElement[1] == "desc" {
				sortParams = append(sortParams, bson.E{sortElement[0], -1})
			}
		}
	}

	if len(sortParams) > 0 {
		sortParams = append(sortParams, bson.E{"_id", -1})
	}

	return sortParams
}

func RequestBodyParser(c *fiber.Ctx, object interface{}) error {
	if err := c.BodyParser(&object); err != nil {
		return err
	}

	if err := validate.Struct(object); err != nil {
		return err
	}

	return nil
}

func RequestBodyParserNoValidate(c *fiber.Ctx, object interface{}) error {
	if err := c.BodyParser(&object); err != nil {
		return err
	}

	return nil
}

func GetStringFromRequestParam(c *fiber.Ctx, paramName string) string {
	return c.Query(paramName, "")
}

func GetBoolFromRequestParam(c *fiber.Ctx, paramName string) (*bool, error) {
	var output *bool = nil
	str := c.Query(paramName, "")

	if str != "" {
		boolValue, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}

		output = &boolValue
	}

	return output, nil
}

func GetObjectIdFromRequestParam(c *fiber.Ctx, paramName string, allowNull bool) (*primitive.ObjectID, error) {
	var output *primitive.ObjectID = nil
	str := c.Query(paramName, "")

	if str != "" || !allowNull {
		objectID, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil, err
		}

		output = &objectID
	}

	return output, nil
}

func GetInt64FromRequestParam(c *fiber.Ctx, paramName string, allowNull bool) (*int64, error) {
	var output *int64 = nil
	str := c.Query(paramName, "")

	if str != "" || !allowNull {
		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}

		output = &value
	}

	return output, nil
}

// string format "yyyyMMdd" in GMT+7 to time
func GetTimeFromRequestParam_Date01_GmtPlus7(c *fiber.Ctx, paramName string, allowNull bool) (*time.Time, error) {
	var output *time.Time = nil
	str := c.Query(paramName, "")

	if str != "" || !allowNull {
		timeValue, err := StringDate01_GmtPlus7_ToTime(str)
		if err != nil {
			return nil, err
		}

		output = &timeValue
	}

	return output, nil
}

func GetObjectIdSliceFromRequestParam(c *fiber.Ctx, paramName string) ([]primitive.ObjectID, error) {
	output := []primitive.ObjectID{}
	queryParams := c.Context().QueryArgs()
	aParams := []string{}

	queryParams.VisitAll(func(key, value []byte) {
		if string(key) == paramName {
			valueStr := string(value)
			if valueStr != "" {
				aParams = append(aParams, string(value))
			}
		}
	})

	for _, v := range aParams {
		if objectID, err := primitive.ObjectIDFromHex(v); err != nil {
			return output, err
		} else {
			output = append(output, objectID)
		}
	}

	return output, nil
}

func GetObjectIdFromRequestPath(c *fiber.Ctx, paramName string) (primitive.ObjectID, error) {
	str := c.Params(paramName)

	if _id, err := primitive.ObjectIDFromHex(str); err != nil {
		return primitive.NilObjectID, err
	} else {
		return _id, nil
	}
}

func GetStringFromRequestPath(c *fiber.Ctx, paramName string) string {
	str := c.Params(paramName)

	if decoded, err := url.QueryUnescape(str); err != nil {
		return str
	} else {
		return decoded
	}
}
