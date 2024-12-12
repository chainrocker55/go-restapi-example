package models

type AccountMeReq struct {
	CustomerId    string
	CisUid        string
	CorrelationId string
}

type AccountMeResp struct {
	Uid      string     `json:"uid"`
	Accounts []Accounts `json:"accounts"`
}

type Accounts struct {
	AccountNo   string      `json:"accountNumber"`
	ProductType string      `json:"productType"`
	AccountType AccountType `json:"type,omitempty"`
}

type AccountType struct {
	Type string `json:"type"`
	Code string `json:"code"`
}
