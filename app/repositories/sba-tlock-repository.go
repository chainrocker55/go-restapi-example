package repositories

import (
	"creditlimit-connector/app/database"
	"sync"

	"gorm.io/gorm"
)

var (
	onceSBATlockRepo sync.Once
	sbaTlockRepo     SBATlockRepo
)

type SBATlockRepo interface {
	FilterValidAccount(accounts []string) ([]string, error)
}

type SBATlockRepoImp struct {
	db *gorm.DB
}

func InitSBATlockRepo() SBATlockRepo {
	onceSBATlockRepo.Do(func() {
		sbaTlockRepo = &SBATlockRepoImp{
			db: database.InitINVXCreditLimitDatabase(),
		}
	})
	return sbaTlockRepo
}


func (r *SBATlockRepoImp) FilterValidAccount(accounts []string) ([]string, error) {
	var data []string
	result := r.db.Where(`SELECT account_no as invalidAccount
						FROM sba_tlock
						WHERE account_no in (?) and xchgmkt = '1'
						and current_date() between effdate and enddate 
						and reason_code in ('LOMB', 'PPBL');`, accounts).Pluck("invalidAccount", &data)
	return data, result.Error
}
