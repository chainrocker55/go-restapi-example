package services

import (
	"context"
	"creditlimit-connector/app/clients"
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/consts"
	"creditlimit-connector/app/entity"
	"creditlimit-connector/app/models"
	"creditlimit-connector/app/repositories"
	"creditlimit-connector/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"creditlimit-connector/app/log"

	"github.com/redis/go-redis/v9"
)

var (
	oncecreditService sync.Once
	creditService     CreditService
)

type CreditService interface {
	CheckCreditLimit(ctx context.Context, product string, category, accountType string) models.CheckCreditLimitResp
	GenerateRefID(runningNo int, date string, suffix string, key string) string
	GetRunningNo(key string) int
	// QueryTempCreditLimit(req models.QueryTempCreditLimitReq) models.QueryTempCreditLimitResp
	// QueryGoodAsset(req models.QueryGoodAssetReq) models.QueryGoodAssetResp
	// AdjustCreditLimit(req models.AdjustCreditLimitReq) models.AdjustCreditLimitResp
	CheckOperationalDate() bool
}

type CreditServiceImp struct {
	creditClient                 clients.CreditClient
	accountClient                clients.AccountClient
	runningRepo                  repositories.RunningNoRepo
	redisRepo                    repositories.RedisRepository
	fnGetOnshoreBusinessDateRepo repositories.FnGetOnshoreBusinessDateRepo
}

func InitCreditService() CreditService {
	oncecreditService.Do(func() {
		creditService = &CreditServiceImp{
			creditClient:                 clients.InitCreditClient(),
			runningRepo:                  repositories.InitRunningNoRepo(),
			redisRepo:                    repositories.InitRedisRepository(),
			fnGetOnshoreBusinessDateRepo: repositories.InitFnGetOnshoreBusinessDateRepo(),
		}
	})
	return creditService
}

func (s *CreditServiceImp) CheckCreditLimit(ctx context.Context, product string, category, accountType string) models.CheckCreditLimitResp {
	panic("implement me")
}

func (s *CreditServiceImp) tranformQueryTempCreditLimit(req models.QueryTempCreditLimitReq) models.QueryTempCreditLimitResp {
	sbaConfig := configs.Conf.SBA
	encryptReq := utils.EncryptSBARequest(req, sbaConfig)
	res := s.creditClient.QueryTempCreditLimit(models.EncryptQueryCreditLimit{Msg: encryptReq})
	decryptRes := decryptResponse[models.QueryTempCreditLimitResp](res)
	return decryptRes
}

func (s *CreditServiceImp) tranformAdjustCreditLimit(req models.AdjustCreditLimitReq) models.AdjustCreditLimitResp {
	sbaConfig := configs.Conf.SBA
	encryptReq := utils.EncryptSBARequest(req, sbaConfig)
	res := s.creditClient.AdjustCreditLimit(models.EncryptQueryCreditLimit{Msg: encryptReq})
	decryptRes := decryptResponse[models.AdjustCreditLimitResp](res)
	return decryptRes
}

// QueryGoodAsset implements CreditService.
func (s *CreditServiceImp) tranformQueryGoodAsset(req models.QueryGoodAssetReq) models.QueryGoodAssetResp {
	sbaConfig := configs.Conf.SBA
	encryptReq := utils.EncryptSBARequest(req, sbaConfig)
	res := s.creditClient.QueryGoodAsset(models.EncryptQueryCreditLimit{Msg: encryptReq})
	decryptRes := decryptResponse[models.QueryGoodAssetResp](res)
	return decryptRes
}

func decryptResponse[R any](res models.EncryptQueryCreditLimit) R {
	sbaConfig := configs.Conf.SBA
	decryptRes := utils.DecryptSBARequest(res, sbaConfig)
	var m R
	json.Unmarshal(decryptRes, &m)
	return m
}

func (s *CreditServiceImp) GenerateRefID(runningNo int, date string, suffix string, key string) string {
	random4 := rand.Intn(10000)
	return configs.Conf.PrefixRefID + date + fmt.Sprintf("%06d", runningNo) + fmt.Sprintf("%04d", random4) + suffix
}

func (s *CreditServiceImp) GetRunningNo(key string) int {
	// Truncate the time part to set it to 00:00:00
	datetime := time.Now().Truncate(24 * time.Hour)
	log.Infof("Finding running no name: %s time: %v", key, datetime)
	runningNo, err := s.runningRepo.FindByNameAndUpdatedAt(key, datetime)
	if err != nil {
		panic(errors.New("error database"))
	}
	if runningNo == nil {
		log.Infof("Not found key: %s", key)
		runningNo = &entity.RunningNoEntity{
			Name:  key,
			Value: 1,
		}
		s.runningRepo.Save(*runningNo)
		return 1
	}

	return runningNo.Value
}

func (s *CreditServiceImp) CheckOperationalDate() bool {
	log.Info("Checking operational date")
	data := getOnshoreBusinessDateFromRedis(s)
	log.Infof("Operational date from redis: %+v", data)
	if data != nil {
		return checkBusinessDate(*data)
	} else {
		data := getOnshoreBusinessDateFromDatabase(s)
		log.Infof("Operational date from database: %+v", data)
		ttlInSeconds := 24 * 60 * 60 // 24 hours
		err := s.redisRepo.Save(consts.OPERATIONAL_DATE_KEY, data, ttlInSeconds)
		log.Info("Operational date saved to redis")
		if err != nil {
			panic("Cannot save operational date to redis")
		}
		return checkBusinessDate(data)
	}
}

func checkBusinessDate(data models.FnGetOnshoreBusinessDateModel) bool {
	return data.IsOperationalDate && isBusinessDateInTimeRange()
}

func isBusinessDateInTimeRange() bool {
	// Get the current time
	currentTime := time.Now()
	startTimeStr := configs.Conf.SBA.OperationTimeStart
	endTimeStr := configs.Conf.SBA.OperationTimeEnd

	// Define the start and end time for the range (9:00 AM and 4:55 PM)
	startTime, _ := time.Parse(time.TimeOnly, startTimeStr)
	endTime, _ := time.Parse(time.TimeOnly, endTimeStr)

	startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime.Hour(), startTime.Minute(), 0, 0, currentTime.Location())
	endTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endTime.Hour(), endTime.Minute(), 0, 0, currentTime.Location())

	log.Infof("Current time: %v, Start time: %v, End time: %v", currentTime.Format(time.TimeOnly), startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly))

	// Check if the current time is within the range
	return currentTime.After(startTime) && currentTime.Before(endTime)
}

func getOnshoreBusinessDateFromRedis(s *CreditServiceImp) *models.FnGetOnshoreBusinessDateModel {
	var model models.FnGetOnshoreBusinessDateModel
	err := s.redisRepo.Find(consts.OPERATIONAL_DATE_KEY, &model)
	if err == redis.Nil {
		log.Info("Operational date not found in redis")
		return nil // Return nil if the key is not found, as no value is available
	}
	if err != nil {
		log.Error(err)
		err = models.NewErrorResponse(422, "04001", "The targeted Redis Cache cannot be connected")
		panic(err)
	}
	return &model
}

func getOnshoreBusinessDateFromDatabase(s *CreditServiceImp) models.FnGetOnshoreBusinessDateModel {
	result, err := s.fnGetOnshoreBusinessDateRepo.FindLocalToday()
	if err != nil {
		log.Error(err)
		err := models.NewErrorResponse(422, "03001", "The targeted database cannot be connected")
		panic(err)
	}
	return result
}
