package qlang

// import "fmt"

// type FILTER_OP_TYPE string
// type _TRANSLATION_TABLE map[FILTER_OP_TYPE]string

// const (
// 	FOP_OR  FILTER_OP_TYPE = "OR"
// 	FOP_AND FILTER_OP_TYPE = "AND"
// 	FOP_EQ  FILTER_OP_TYPE = "EQ"
// 	FOP_GT  FILTER_OP_TYPE = "GT"
// 	FOP_GTE FILTER_OP_TYPE = "GTE"
// 	FOP_LT  FILTER_OP_TYPE = "LT"
// 	FOP_LTE FILTER_OP_TYPE = "LTE"
// 	FOP_LE  FILTER_OP_TYPE = "LE"
// )

// var fop_priorities = map[FILTER_OP_TYPE]int{
// 	FOP_OR:  3,
// 	FOP_AND: 2,
// 	FOP_EQ:  1,
// 	FOP_GT:  1,
// 	FOP_GTE: 1,
// 	FOP_LT:  1,
// 	FOP_LTE: 1,
// 	FOP_LE:  1,
// }

// var FopTable = map[string]FILTER_OP_TYPE{
// 	"OR":  FOP_OR,
// 	"AND": FOP_AND,
// 	"EQ":  FOP_EQ,
// 	"GT":  FOP_EQ,
// 	"GTE": FOP_GTE,
// 	"LT":  FOP_LT,
// 	"LTE": FOP_LTE,
// 	"LE":  FOP_LE,
// }

// type FilterOpItem struct {
// 	OperationType FILTER_OP_TYPE
// 	Priority      int
// 	// tTable        _TRANSLATION_TABLE
// }

// func NewFilterOpItem(fopt FILTER_OP_TYPE) (*FilterOpItem, error) {
// 	p, ok := getPriority(fopt)
// 	if !ok {
// 		return nil, fmt.Errorf("%s operator does not have a priority", fopt)
// 	}

// 	return &FilterOpItem{OperationType: fopt, Priority: p}, nil
// }

// func getPriority(fopt FILTER_OP_TYPE) (int, bool) {
// 	val, ok := fop_priorities[fopt]
// 	return val, ok
// }

// // func IsOperator
