package repositories

import (
	"creditlimit-connector/app/database"
	"creditlimit-connector/app/entity"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	onceRunningNoRepo sync.Once
	runningNoRepo     RunningNoRepo
)

type RunningNoRepo interface {
	Save(data entity.RunningNoEntity) error
	FindByNameAndUpdatedAt(name string, datetime time.Time) (*entity.RunningNoEntity, error)
}

type RunningNoRepoImp struct {
	db *gorm.DB
}

func InitRunningNoRepo() RunningNoRepo {
	onceRunningNoRepo.Do(func() {
		runningNoRepo = &RunningNoRepoImp{
			db: database.InitINVXCreditLimitDatabase(),
		}
	})
	return runningNoRepo
}

func (r *RunningNoRepoImp) Save(data entity.RunningNoEntity) error {
	result := r.db.Save(&data)
	return result.Error
}

func (r *RunningNoRepoImp) FindByNameAndUpdatedAt(name string, datetime time.Time) (*entity.RunningNoEntity, error) {
	var runningNo entity.RunningNoEntity
	result := r.db.Where("name = ? and DATE(updated_at) >= ?", name, datetime).First(&runningNo)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &runningNo, result.Error
}
