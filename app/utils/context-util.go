package utils

import (
	"creditlimit-connector/app/consts"
	"creditlimit-connector/app/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/modern-go/gls"
)

func GetContextValue(key string) string {
	value := gls.Get(key)
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return ""
}

func GenerateID() string {
	return uuid.New().String()
}

func GetRequestID() string {
	return GetContextValue(consts.ContextRequestId)
}

func GetCorrelationId() string {
	return GetContextValue(consts.ContextCorrelationId)
}

func GetOriginatingIp() string {
	return GetContextValue(consts.ContextOriginatingIp)
}

func GetCustomerId() string {
	return GetContextValue(consts.ContextCustomerId)
}

func GetCisUid() string {
	return GetContextValue(consts.ContextCisUid)
}

func GetRequiredParam(key string, c *fiber.Ctx) string {
	value := c.Params(key)
	if value == "" {
		message := fmt.Sprintf("Missing required parameter: %s", key)
		panic(models.NewErrorResponse(406, "461", message))
	}
	return value
}
