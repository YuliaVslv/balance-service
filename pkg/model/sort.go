package model

type SortField string

const (
	None    SortField = "id"
	ByValue           = "value"
	ByDate            = "date"
)

type Order string

const (
	ASC  Order = "ASC"
	DESC       = "DESC"
)
