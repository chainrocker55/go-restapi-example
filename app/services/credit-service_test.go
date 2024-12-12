package services

import (
	"context"
	mocksClient "creditlimit-connector/app/clients/mocks"
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/consts"
	"creditlimit-connector/app/entity"
	"creditlimit-connector/app/models"
	mocksRepositories "creditlimit-connector/app/repositories/mocks"
	"creditlimit-connector/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckCreditLimit(t *testing.T) {
	mockCreditClient := mocksClient.NewCreditClient(t)
	mockRepo := mocksRepositories.NewRunningNoRepo(t)

	service := &CreditServiceImp{
		creditClient: mockCreditClient,
		runningRepo:  mockRepo,
	}

	ctx := context.Background()
	accountType := "1234567890"
	product := "testProduct"
	category := "testCategory"

	t.Run("Success", func(t *testing.T) {

		resp := service.CheckCreditLimit(ctx, product, category, accountType)
		assert.NotNil(t, resp)
	})
}

func TestGenerateRefID(t *testing.T) {
	service := &CreditServiceImp{}

	tests := []struct {
		name        string
		runningNo   int
		date        string
		suffix      string
		key         string
		prefixRefID string
	}{
		{
			name:        "Generate RefID with valid inputs",
			runningNo:   123,
			date:        "20230101",
			suffix:      "XYZ",
			key:         "testKey",
			prefixRefID: "PREFIX",
		},
		{
			name:        "Generate RefID with different inputs",
			runningNo:   456,
			date:        "20230202",
			suffix:      "ABC",
			key:         "anotherKey",
			prefixRefID: "PREFIX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configs.Conf.PrefixRefID = tt.prefixRefID
			refID := service.GenerateRefID(tt.runningNo, tt.date, tt.suffix, tt.key)

			assert.Contains(t, refID, tt.prefixRefID)
			assert.Contains(t, refID, tt.date)
			assert.Contains(t, refID, fmt.Sprintf("%06d", tt.runningNo))
			assert.Contains(t, refID, tt.suffix)
			assert.Len(t, refID, len(tt.prefixRefID)+len(tt.date)+6+4+len(tt.suffix))
		})
	}
}

func TestGetRunningNo(t *testing.T) {

	key := "testKey"
	datetime := time.Now().Truncate(24 * time.Hour)

	t.Run("Record found", func(t *testing.T) {
		mockRepo := mocksRepositories.NewRunningNoRepo(t)
		service := &CreditServiceImp{
			runningRepo: mockRepo,
		}
		mockRepo.On("FindByNameAndUpdatedAt", mock.Anything, mock.Anything).Return(&entity.RunningNoEntity{Value: 5}, nil)

		runningNo := service.GetRunningNo(key)

		assert.Equal(t, 5, runningNo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Record not found", func(t *testing.T) {
		mockRepo := mocksRepositories.NewRunningNoRepo(t)
		service := &CreditServiceImp{
			runningRepo: mockRepo,
		}
		mockRepo.On("FindByNameAndUpdatedAt", key, datetime).Return(nil, nil)
		mockRepo.On("Save", mock.Anything).Return(nil)

		runningNo := service.GetRunningNo(key)

		assert.Equal(t, 1, runningNo)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		mockRepo := mocksRepositories.NewRunningNoRepo(t)
		service := &CreditServiceImp{
			runningRepo: mockRepo,
		}
		mockRepo.On("FindByNameAndUpdatedAt", key, datetime).Return(nil, errors.New("database error"))

		assert.Panics(t, func() {
			service.GetRunningNo(key)
		})

		mockRepo.AssertExpectations(t)
	})
}

func TestQueryTempCreditLimit(t *testing.T) {

	req := models.QueryTempCreditLimitReq{
		AuthorizedKey: "testAuthKey",
		ReferID:       "testRefID",
		SenderID:      "testSenderID",
		AccountNo:     "testAccountNo",
		SendDate:      "20230101",
		SendTime:      "120000",
	}

	t.Run("Success", func(t *testing.T) {
		mockCreditClient := mocksClient.NewCreditClient(t)
		key := "1234567890123456"

		configs.Conf.SBA.MessageKey = key

		service := &CreditServiceImp{
			creditClient: mockCreditClient,
		}

		expected := models.QueryTempCreditLimitResp{
			ReferID:    "testRefID",
			ResultCode: "0000",
		}

		jsonData, _ := json.Marshal(expected)
		encryptedResBody, _ := utils.EncryptAES128ECB(jsonData, key)
		resClient := models.EncryptQueryCreditLimit{Msg: encryptedResBody}

		mockCreditClient.On("QueryTempCreditLimit", mock.Anything).Return(resClient)

		resp := service.queryTempCreditLimit(req)

		mockCreditClient.AssertExpectations(t)
		assert.Equal(t, "testRefID", resp.ReferID)
		assert.Equal(t, "0000", resp.ResultCode)
	})
}

func TestAdjustCreditLimit(t *testing.T) {

	req := models.AdjustCreditLimitReq{
		AuthorizedKey: "testAuthKey",
		ReferID:       "testRefID",
		SenderID:      "testSenderID",
		AccountNo:     "testAccountNo",
		SendDate:      "20230101",
		SendTime:      "120000",
	}

	t.Run("Success", func(t *testing.T) {
		mockCreditClient := mocksClient.NewCreditClient(t)
		key := "1234567890123456"

		configs.Conf.SBA.MessageKey = key

		service := &CreditServiceImp{
			creditClient: mockCreditClient,
		}

		expected := models.AdjustCreditLimitResp{
			ReferID:    "testRefID",
			ResultCode: "0000",
		}

		jsonData, _ := json.Marshal(expected)
		encryptedResBody, _ := utils.EncryptAES128ECB(jsonData, key)
		resClient := models.EncryptQueryCreditLimit{Msg: encryptedResBody}

		mockCreditClient.On("AdjustCreditLimit", mock.Anything).Return(resClient)

		resp := service.adjustCreditLimit(req)

		mockCreditClient.AssertExpectations(t)
		assert.Equal(t, "testRefID", resp.ReferID)
		assert.Equal(t, "0000", resp.ResultCode)
	})
}

func TestQueryGoodAsset(t *testing.T) {

	req := models.QueryGoodAssetReq{
		AuthorizedKey: "testAuthKey",
		ReferID:       "testRefID",
		SenderID:      "testSenderID",
		SendDate:      "20230101",
		SendTime:      "120000",
	}

	t.Run("Success", func(t *testing.T) {
		mockCreditClient := mocksClient.NewCreditClient(t)
		key := "1234567890123456"

		configs.Conf.SBA.MessageKey = key

		service := &CreditServiceImp{
			creditClient: mockCreditClient,
		}

		expected := models.QueryGoodAssetResp{
			ReferID:    "testRefID",
			ResultCode: "0000",
		}

		jsonData, _ := json.Marshal(expected)
		encryptedResBody, _ := utils.EncryptAES128ECB(jsonData, key)
		resClient := models.EncryptQueryCreditLimit{Msg: encryptedResBody}

		mockCreditClient.On("QueryGoodAsset", mock.Anything).Return(resClient)

		resp := service.queryGoodAsset(req)

		mockCreditClient.AssertExpectations(t)
		assert.Equal(t, "testRefID", resp.ReferID)
		assert.Equal(t, "0000", resp.ResultCode)
	})
}

func TestCheckOperationalDate(t *testing.T) {
configs.Conf.SBA.OperationTimeStart = "00:00:00"
		configs.Conf.SBA.OperationTimeEnd = "23:59:59"
	t.Run("Success from redis", func(t *testing.T) {
		// Arrange

		
		mockRedisRepo := mocksRepositories.NewRedisRepository(t)
		mockFnGetOnshoreBusinessDateRepo := mocksRepositories.NewFnGetOnshoreBusinessDateRepo(t)
		expectedModel := &models.FnGetOnshoreBusinessDateModel{
			Date: 			"2021-01-01",
			IsOperationalDate: true,
		}

		// Mock the behavior of redisRepo.Find
		mockRedisRepo.On("Find", consts.OPERATIONAL_DATE_KEY, mock.Anything).Run(func(args mock.Arguments) {
			model := args.Get(1).(*models.FnGetOnshoreBusinessDateModel)
			// Set up model with expected values
			*model = *expectedModel
		}).Return(nil)

		// fnGetOnshoreBusinessDateRepo.On("FindLocalToday").Return(expectedModel)

		// Create CreditServiceImp with the mock
		service := &CreditServiceImp{
			redisRepo: mockRedisRepo,
			fnGetOnshoreBusinessDateRepo: mockFnGetOnshoreBusinessDateRepo,
		}

		// Act
		result := service.CheckOperationalDate()

		// Assert
		assert.NotNil(t, result)
		assert.Equal(t, true, result)
		mockRedisRepo.AssertExpectations(t)
		mockRedisRepo.AssertNotCalled(t, "Save")
		mockFnGetOnshoreBusinessDateRepo.AssertNotCalled(t, "FindLocalToday")
	})

	t.Run("Success from database", func(t *testing.T) {
		// Arrange
		mockRedisRepo := mocksRepositories.NewRedisRepository(t)
		mockFnGetOnshoreBusinessDateRepo := mocksRepositories.NewFnGetOnshoreBusinessDateRepo(t)
		expectedModel := models.FnGetOnshoreBusinessDateModel{
			Date: 			"2021-01-01",
			IsOperationalDate: true,
		}

		// Mock the behavior of redisRepo.Find
		mockRedisRepo.On("Find", consts.OPERATIONAL_DATE_KEY, mock.Anything).Return(redis.Nil)
		mockRedisRepo.On("Save", consts.OPERATIONAL_DATE_KEY, mock.Anything, mock.Anything).Return(nil)
		mockFnGetOnshoreBusinessDateRepo.On("FindLocalToday").Return(expectedModel, nil)

		// Create CreditServiceImp with the mock
		service := &CreditServiceImp{
			redisRepo: mockRedisRepo,
			fnGetOnshoreBusinessDateRepo: mockFnGetOnshoreBusinessDateRepo,
		}

		// Act
		result := service.CheckOperationalDate()

		// Assert
		assert.NotNil(t, result)
		assert.Equal(t, true, result)
		mockRedisRepo.AssertExpectations(t)
		mockRedisRepo.AssertCalled(t, "Find", consts.OPERATIONAL_DATE_KEY, mock.Anything)
		mockRedisRepo.AssertCalled(t, "Save", consts.OPERATIONAL_DATE_KEY, mock.Anything, mock.Anything)
		mockFnGetOnshoreBusinessDateRepo.AssertCalled(t, "FindLocalToday")
	})

	t.Run("Failed from redis", func(t *testing.T) {
		// Arrange
		mockRedisRepo := mocksRepositories.NewRedisRepository(t)
		mockFnGetOnshoreBusinessDateRepo := mocksRepositories.NewFnGetOnshoreBusinessDateRepo(t)

		// Mock the behavior of redisRepo.Find
		mockRedisRepo.On("Find", consts.OPERATIONAL_DATE_KEY, mock.Anything).Return(errors.New("error"))
		// Create CreditServiceImp with the mock
		service := &CreditServiceImp{
			redisRepo: mockRedisRepo,
			fnGetOnshoreBusinessDateRepo: mockFnGetOnshoreBusinessDateRepo,
		}

		// Act
		assert.Panics(t, func() {
			service.CheckOperationalDate()
		})

		// Assert
		mockRedisRepo.AssertCalled(t, "Find", consts.OPERATIONAL_DATE_KEY, mock.Anything)
		mockRedisRepo.AssertNotCalled(t, "Save", consts.OPERATIONAL_DATE_KEY, mock.Anything, mock.Anything)
		mockFnGetOnshoreBusinessDateRepo.AssertNotCalled(t, "FindLocalToday")
	})
}

