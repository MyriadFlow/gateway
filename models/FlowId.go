package models

type FlowIdType string

const (
	AUTH FlowIdType = "AUTH"
)

type FlowId struct {
	FlowId        string     `gorm:"primary_key"`
	FlowIdType    FlowIdType `json:"flow_id_type"`
	WalletAddress string     `json:"wallet_address"`
}
