package model

type GraviteeAMListWrapper struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"currentPage"`
	TotalCount  int         `json:"totalCount"`
}
