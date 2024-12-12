package clients

import (
	"creditlimit-connector/app/models"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryTempCreditLimit(t *testing.T) {
	client := resty.New()
	service := &CreditClientImp{
		client: client,
	}

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("Case_Network_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}

		mockErr := errors.New("Something when wrong")
		httpmock.RegisterResponder("POST", `=~.*/QueryTempCreditLimit$`,
			httpmock.NewErrorResponder(mockErr))

		assert.Panics(t, func() {
			service.QueryTempCreditLimit(mockRequest)
		})

	})

	t.Run("Case_Unmarshal_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		expected := `"orderNo": "100.0"`
		httpmock.RegisterResponder("POST", `=~.*/QueryTempCreditLimit$`,
			httpmock.NewStringResponder(200, expected))

		assert.Panics(t, func() {
			service.QueryTempCreditLimit(mockRequest)
		})

	})

	t.Run("Case_400", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		httpmock.RegisterResponder("POST", `=~.*/QueryTempCreditLimit$`,
			httpmock.NewStringResponder(400, ""))

		assert.Panics(t, func() {
			service.QueryTempCreditLimit(mockRequest)
		})
	})

	t.Run("Case_Success", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test-request",
		}

		expected := models.EncryptQueryCreditLimit{
			Msg: "test-response",
		}

		jsonData, _ := json.Marshal(expected)

		httpmock.RegisterResponder("POST", `=~.*/QueryTempCreditLimit$`,
			httpmock.NewBytesResponder(200, jsonData))

		res := service.QueryTempCreditLimit(mockRequest)
		assert.Equal(t, expected, res)
	})

}


func TestQueryGoodAsset(t *testing.T) {
	client := resty.New()
	service := &CreditClientImp{
		client: client,
	}

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("Case_Network_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}

		mockErr := errors.New("Something when wrong")
		httpmock.RegisterResponder("POST", `=~.*/QueryGoodAsset$`,
			httpmock.NewErrorResponder(mockErr))

		assert.Panics(t, func() {
			service.QueryGoodAsset(mockRequest)
		})

	})

	t.Run("Case_Unmarshal_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		expected := `"orderNo": "100.0"`
		httpmock.RegisterResponder("POST", `=~.*/QueryGoodAsset$`,
			httpmock.NewStringResponder(200, expected))

		assert.Panics(t, func() {
			service.QueryGoodAsset(mockRequest)
		})

	})

	t.Run("Case_400", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		httpmock.RegisterResponder("POST", `=~.*/QueryGoodAsset$`,
			httpmock.NewStringResponder(400, ""))

		assert.Panics(t, func() {
			service.QueryGoodAsset(mockRequest)
		})
	})

	t.Run("Case_Success", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test-request",
		}

		expected := models.EncryptQueryCreditLimit{
			Msg: "test-response",
		}

		jsonData, _ := json.Marshal(expected)

		httpmock.RegisterResponder("POST", `=~.*/QueryGoodAsset$`,
			httpmock.NewBytesResponder(200, jsonData))

		res := service.QueryGoodAsset(mockRequest)
		assert.Equal(t, expected, res)
	})

}

func TestAdjustCreditLimit(t *testing.T) {
	client := resty.New()
	service := &CreditClientImp{
		client: client,
	}

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	t.Run("Case_Network_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}

		mockErr := errors.New("Something when wrong")
		httpmock.RegisterResponder("POST", `=~.*/AdjustCreditLimit$`,
			httpmock.NewErrorResponder(mockErr))

		assert.Panics(t, func() {
			service.AdjustCreditLimit(mockRequest)
		})

	})

	t.Run("Case_Unmarshal_error", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		expected := `"orderNo": "100.0"`
		httpmock.RegisterResponder("POST", `=~.*/AdjustCreditLimit$`,
			httpmock.NewStringResponder(200, expected))

		assert.Panics(t, func() {
			service.AdjustCreditLimit(mockRequest)
		})

	})

	t.Run("Case_400", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test",
		}
		httpmock.RegisterResponder("POST", `=~.*/AdjustCreditLimit$`,
			httpmock.NewStringResponder(400, ""))

		assert.Panics(t, func() {
			service.AdjustCreditLimit(mockRequest)
		})
	})

	t.Run("Case_Success", func(t *testing.T) {
		mockRequest := models.EncryptQueryCreditLimit{
			Msg: "test-request",
		}

		expected := models.EncryptQueryCreditLimit{
			Msg: "test-response",
		}

		jsonData, _ := json.Marshal(expected)

		httpmock.RegisterResponder("POST", `=~.*/AdjustCreditLimit$`,
			httpmock.NewBytesResponder(200, jsonData))

		res := service.AdjustCreditLimit(mockRequest)
		assert.Equal(t, expected, res)
	})

}

