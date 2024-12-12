package controllers

import (
	"creditlimit-connector/app/consts"
	"creditlimit-connector/app/log"
	"creditlimit-connector/app/models"
	"creditlimit-connector/app/services"
	"creditlimit-connector/app/utils"

	"github.com/gofiber/fiber/v2"
)

type CreditController interface {
	CheckCreditLimit(c *fiber.Ctx) error
	AdjustCreditLimit(c *fiber.Ctx) error
	CheckOperationalDate(c *fiber.Ctx) error
}

type CreditControllerImp struct {
	creditService services.CreditService
}

func InitCreditController() CreditController {
	creditService := services.InitCreditService()
	return &CreditControllerImp{
		creditService: creditService,
	}
}

func (s *CreditControllerImp) CheckCreditLimit(c *fiber.Ctx) error {
	product := utils.GetRequiredParam("product", c)
	accountType := utils.GetRequiredParam("accountType", c)
	category := utils.GetRequiredParam("category", c)

	log.Infof("CheckCreditLimit: product=%s, accountType=%s, category=%s", product, accountType, category)

	if product != consts.EQUITY {
		return c.JSON(models.NewErrorResponse(406, "463", "The field product does not match the values"))
	}

	if category != consts.ONSHORE {
		return c.JSON(models.NewErrorResponse(406, "463", "The field category does not match the values"))
	}

	if accountType != consts.CASH_BALANCE {
		return c.JSON(models.NewErrorResponse(406, "463", "The field accountType does not match the values"))
	}

	s.creditService.CheckCreditLimit(c.Context(), product, category, accountType)
	return c.JSON("")
}

func (s *CreditControllerImp) CheckOperationalDate(c *fiber.Ctx) error {
	isOptionalDate := s.creditService.CheckOperationalDate()
	return c.JSON(isOptionalDate)
}

// AdjustCreditLimit implements CreditController.
func (s *CreditControllerImp) AdjustCreditLimit(c *fiber.Ctx) error {
	panic("unimplemented")
}
