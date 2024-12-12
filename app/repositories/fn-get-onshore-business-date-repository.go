package repositories

import (
	"creditlimit-connector/app/database"
	"creditlimit-connector/app/models"
	"sync"

	"gorm.io/gorm"
)

var (
	onceFnGetOnshoreBusinessDate sync.Once
	fnGetOnshoreBusinessDate     FnGetOnshoreBusinessDateRepo
)

type FnGetOnshoreBusinessDateRepo interface {
	FindLocalToday() (models.FnGetOnshoreBusinessDateModel, error)
}

type FnGetOnshoreBusinessDateRepoImp struct {
	db *gorm.DB
}

func InitFnGetOnshoreBusinessDateRepo() FnGetOnshoreBusinessDateRepo {
	onceFnGetOnshoreBusinessDate.Do(func() {
		fnGetOnshoreBusinessDate = &FnGetOnshoreBusinessDateRepoImp{
			db: database.InitINVXCenterDatabase(),
		}
	})
	return fnGetOnshoreBusinessDate
}

func (r *FnGetOnshoreBusinessDateRepoImp) FindLocalToday() (models.FnGetOnshoreBusinessDateModel, error) {
	var data models.FnGetOnshoreBusinessDateModel
	sql := `SELECT dat1.local_today as date,
	   CASE WHEN DATENAME(WEEKDAY, dat1.local_today) in ('Saturday', 'Sunday') OR doh.holiday_date IS NOT NULL THEN 'False'
	        ELSE 'True' END as is_operational_date 
		FROM ( SELECT CAST( GETDATE() AT TIME ZONE 'UTC' AT TIME ZONE 'SE Asia Standard Time' as DATE ) as local_today ) dat1 
	 	left join silver.DML_ONSHORE_HOLIDAY doh on dat1.local_today= doh.holiday_date;`
	result := r.db.Raw(sql).Scan(&data)

	return data, result.Error
}
