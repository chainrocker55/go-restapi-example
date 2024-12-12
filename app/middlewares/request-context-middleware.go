package middlewares

import (
	"creditlimit-connector/app/consts"
	"creditlimit-connector/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/modern-go/gls"
)

func RequestContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(consts.HeaderRequestId)
		if requestID == "" {
			requestID = utils.GenerateID()
		}
		correlationID := c.Get(consts.HeaderCorrelationId)
		if correlationID == "" {
			correlationID = utils.GenerateID()
		}

		originatingIP := c.Get(consts.HeaderOriginatingIp)
		clientIP := c.Get(consts.HeaderClientIp)
		gls.ResetGls(gls.GoID(), make(map[interface{}]interface{}))
		gls.Set(consts.ContextRequestId, requestID)
		gls.Set(consts.ContextCorrelationId, correlationID)
		gls.Set(consts.ContextOriginatingIp, utils.If[string](originatingIP != "", originatingIP, clientIP))
		gls.Set(consts.ContextCustomerId, c.Get(consts.HeaderCustomerId))
		gls.Set(consts.ContextCisUid, c.Get(consts.HeaderCisUid))
		return c.Next()
	}
}
