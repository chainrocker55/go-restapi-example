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
	onceCreditClient sync.Once
	creditClient     CreditClient
)

type CreditClient interface {
	QueryTempCreditLimit(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit
	QueryGoodAsset(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit
	AdjustCreditLimit(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit
}

type CreditClientImp struct {
	client *resty.Client
}

func InitCreditClient() CreditClient {
	onceCreditClient.Do(func() {
		client := resty.New()
		baseUrl := configs.Conf.Client.CreditUrl + "/GATEWAYSERVICE/CustomeModule"
		client = client.
			SetTimeout(30*time.Second).
			SetBaseURL(baseUrl).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Language", "TH").
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

		creditClient = &CreditClientImp{
			client: client,
		}
	})
	return creditClient
}

func (s *CreditClientImp) QueryTempCreditLimit(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit {
	resp, err := s.client.R().
		SetBody(req).
		Post("/QueryTempCreditLimit")

	if resp.IsError() {
		err := resp.Error().(*models.ErrorResponse)
		panic(models.NewErrorResponse(resp.StatusCode(), err.Code, err.Message))
	}
	if err != nil {
		handleCreditError(err)
	}
	var model models.EncryptQueryCreditLimit
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return model
}

func (s *CreditClientImp) QueryGoodAsset(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit {
	resp, err := s.client.R().
		SetBody(req).
		Post("/QueryGoodAsset")

	if resp.IsError() {
		err := resp.Error().(*models.ErrorResponse)
		panic(models.NewErrorResponse(resp.StatusCode(), err.Code, err.Message))
	}
	if err != nil {
		handleCreditError(err)
	}
	var model models.EncryptQueryCreditLimit
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return model
}

func (s *CreditClientImp) AdjustCreditLimit(req models.EncryptQueryCreditLimit) models.EncryptQueryCreditLimit {
	resp, err := s.client.R().
		SetBody(req).
		Post("/AdjustCreditLimit")

	if resp.IsError() {
		err := resp.Error().(*models.ErrorResponse)
		panic(models.NewErrorResponse(resp.StatusCode(), err.Code, err.Message))

	}
	if err != nil {
		handleCreditError(err)
	}
	var model models.EncryptQueryCreditLimit
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return model
}

func handleCreditError(err error) {
	log.Error(err)
	err = models.NewErrorResponse(422, "020002", "The targeted API (Credit) is down or out of service and cannot be called.")
	panic(err)
}