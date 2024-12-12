package clients

import (
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/models"
	"encoding/json"
	"sync"
	"time"

	"creditlimit-connector/app/log"

	"github.com/go-resty/resty/v2"
)

var (
	onceAccountClient sync.Once
	acountClient      AccountClient
)

type AccountClient interface {
	AccountMe(req models.AccountMeReq) models.AccountMeResp
}

type AccountClientImp struct {
	client *resty.Client
}

func InitAccountClient() AccountClient {
	onceAccountClient.Do(func() {
		client := resty.New()
		baseUrl := configs.Conf.Client.AccountUrl + "/account"
		client = client.
			SetTimeout(30*time.Second).
			SetBaseURL(baseUrl).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Language", "TH").
			SetHeader("X-Api-Key", configs.Conf.Client.CreditUrl).
			OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
				// Now you have access to Client and current Request object
				bodyStr, _ := json.Marshal(req.Body)
				log.Infof("Request: %s:%s body: %s", req.Method, req.URL, string(bodyStr))

				return nil // if its success otherwise return error
			}).
			OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
				// Now you have access to Client and current Response object
				log.Infof("Response: %s:%s [%d] body: %s", resp.Request.Method, resp.Request.URL, resp.StatusCode(), resp.String())
				return nil // if its success otherwise return error
			})

		acountClient = &AccountClientImp{
			client: client,
		}
	})
	return acountClient
}

func (s *AccountClientImp) AccountMe(req models.AccountMeReq) models.AccountMeResp {
	resp, err := s.client.R().
		SetBody(req).
		SetHeader("x-correlation-id", req.CorrelationId).
		SetHeader("x-customer-id", req.CustomerId).
		SetHeader("x-cis-uid", req.CisUid).
		Get("/account/me")

	if resp.IsError() {
		err := resp.Error().(*models.ErrorResponse)
		panic(models.NewErrorResponse(resp.StatusCode(), err.Code, err.Message))
	}
	if err != nil {
		handleAccountError(err)
	}
	var model models.AccountMeResp
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return model
}

func handleAccountError(err error) {
	log.Error(err)
	err = models.NewErrorResponse(422, "020001", "The targeted API (Account) is down or out of service and cannot be called.")
	panic(err)
}
