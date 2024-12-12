package models

type QueryTempCreditLimitReq struct {
	AuthorizedKey      string `json:"authorized_key"`
	ReferID            string `json:"refer_id"`
	AccountNo          string `json:"account_no"`
	EffectiveDateFrom  string `json:"effective_date_from,omitempty"`
	EffectiveDateTo    string `json:"effective_date_to,omitempty"`
	TempCreditLineFlag rune   `json:"temp_credit_line_flag"`
	GetLatest          rune   `json:"get_latest"`
	ExchangeMarket     rune   `json:"exchange_market"`
	SenderID           string `json:"sender_id"`
	SendDate           string `json:"send_date"`
	SendTime           string `json:"send_time"`
}
type QueryTempCreditLimitResp struct {
	ReferID            string   `json:"refer_id"`
	ResultCode         string   `json:"result_code"`
	Reason             string   `json:"reason"`
	SenderID           string   `json:"sender_id"`
	SendDate           string   `json:"send_date"`
	SendTime           string   `json:"send_time"`
	AccountNo          string   `json:"account_no"`
	EffectiveDateTotal int      `json:"effective_date_total"`
	EffectiveDateList  []string `json:"effective_date_list"`
	EffectiveDate      string   `json:"effective_date"`
	ExchangeMarket     rune     `json:"exchange_market"`
	TempCreditLineFlag rune     `json:"temp_credit_line_flag"`
	RequestEndDate     string   `json:"request_end_date"`
	EndDate            string   `json:"end_date"`
	OriginalCreditLine float64  `json:"original_credit_line"`
	RequestCreditLine  float64  `json:"request_credit_line"`
	ApproveCreditLine  float64  `json:"approve_credit_line"`
}

type QueryGoodAssetReq struct {
	AuthorizedKey string   `json:"authorized_key"`
	ReferID       string   `json:"refer_id"`
	AccountNo     []string `json:"account_no"`
	SenderID      string   `json:"sender_id"`
	SendDate      string   `json:"send_date"`
	SendTime      string   `json:"send_time"`
}
type QueryGoodAssetResp struct {
	ReferID            string   `json:"refer_id"`
	ResultCode         string   `json:"result_code"`
	Reason             string   `json:"reason"`
	SenderID           string   `json:"sender_id"`
	SendDate           string   `json:"send_date"`
	SendTime           string   `json:"send_time"`
	ResultsListTotal   int      `json:"results_list_total"`
	EffectiveDateTotal int      `json:"effective_date_total"`
	ResultsList        []string `json:"results_list"`
	AccountNo          string   `json:"account_no"`
	GoodAsset          float64  `json:"good_asset"`
	Cashbalance        float64  `json:"cashbalance"`
	APAccumulate       float64  `json:"ap_accumulate"`
	ARAccumulate       float64  `json:"ar_accumulate"`
}

type AdjustCreditLimitReq struct {
	AuthorizedKey      string   `json:"authorized_key"`
	ReferID            string   `json:"refer_id"`
	TransID            string   `json:"trans_id"`
	AccountNo          string   `json:"account_no"`
	SenderID           string   `json:"sender_id"`
	SendDate           string   `json:"send_date"`
	SendTime           string   `json:"send_time"`
	AccountsList       []string `json:"accounts_list"`
	ExchangeMarket     rune     `json:"exchange_market"`
	TempCreditLineFlag rune     `json:"temp_credit_line_flag"`
	RequestExpireDate  rune     `json:"request_expire_date"`
	ExpireDate         string   `json:"expire_date"`
	RequestCreditLine  float64  `json:"request_credit_line"`
	AppCreditLine      float64  `json:"app_credit_line"`
	Remark1            string   `json:"remark1"`
	Remark2            string   `json:"remark2"`
}
type AdjustCreditLimitResp struct {
	ReferID          string   `json:"refer_id"`
	TransID          string   `json:"trans_id"`
	ResultCode       string   `json:"result_code"`
	Reason           string   `json:"reason"`
	SenderID         string   `json:"sender_id"`
	SendDate         string   `json:"send_date"`
	SendTime         string   `json:"send_time"`
	ResultsListTotal int      `json:"results_list_total"`
	ResultsList      []string `json:"results_list"`
	AccountNo        string   `json:"account_no"`
}

type EncryptQueryCreditLimit struct {
	Msg string `json:"msg"`
}
