package quang

import (
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DATABASE_TYPE int

const (
	MONGO DATABASE_TYPE = iota
	SQL
)

type _OUTPUT_TYPE int

const (
	BSON _OUTPUT_TYPE = iota
	STR
	NUM
)

type filter_op_type string
type mongo_op_type string

const (
	FOP_OR  filter_op_type = "OR"
	FOP_AND filter_op_type = "AND"
	FOP_EQ  filter_op_type = "EQ"
	FOP_GT  filter_op_type = "GT"
	FOP_GTE filter_op_type = "GTE"
	FOP_LT  filter_op_type = "LT"
	FOP_LTE filter_op_type = "LTE"
	FOP_LE  filter_op_type = "LE"
	MOP_OR  mongo_op_type  = "$or"
	MOP_AND mongo_op_type  = "$and"
	MOP_EQ  mongo_op_type  = "$eq"
	MOP_GT  mongo_op_type  = "$gt"
	MOP_GTE mongo_op_type  = "$gte"
	MOP_LT  mongo_op_type  = "$lt"
	MOP_LTE mongo_op_type  = "$lte"
	MOP_LE  mongo_op_type  = "$regex"
)

var fop_priorities = map[filter_op_type]int{
	FOP_OR:  3,
	FOP_AND: 2,
	FOP_EQ:  1,
	FOP_GT:  1,
	FOP_GTE: 1,
	FOP_LT:  1,
	FOP_LTE: 1,
	FOP_LE:  1,
}

var fop_table = map[string]filter_op_type{
	"OR":  FOP_OR,
	"AND": FOP_AND,
	"EQ":  FOP_EQ,
	"GT":  FOP_EQ,
	"GTE": FOP_GTE,
	"LT":  FOP_LT,
	"LTE": FOP_LTE,
	"LE":  FOP_LE,
}

var fop_translation_tb_mongo = map[filter_op_type]mongo_op_type{
	FOP_OR:  MOP_OR,
	FOP_AND: MOP_AND,
	FOP_EQ:  MOP_EQ,
	FOP_GT:  MOP_GT,
	FOP_GTE: MOP_GTE,
	FOP_LT:  MOP_LT,
	FOP_LTE: MOP_LTE,
	FOP_LE:  MOP_LE,
}

type filter_out_i struct {
	val     any
	valType _OUTPUT_TYPE
}

type filter_op_i struct {
	t filter_op_type
	p int
}

type FilterTranslator struct {
	outputStack   *Stack[filter_out_i]
	operatorStack *Stack[filter_op_i]
}

func NewFilterTranslator() *FilterTranslator {
	return &FilterTranslator{
		outputStack:   NewStack[filter_out_i](),
		operatorStack: NewStack[filter_op_i](),
	}
}

func (ft *FilterTranslator) Translate(s string, t DATABASE_TYPE) (*bson.D, error) {
	switch t {
	case MONGO:
		return ft.TranslateToMongo(s)
	case SQL:
		fallthrough
	default:
		panic("unimplemented")
	}
}

// Implements the shunting yard algorithm for operator-precedence parsing
func (ft *FilterTranslator) TranslateToMongo(s string) (*bson.D, error) {
	_LEN := len(s)

	// Scan from left to right
	for i := 0; i < _LEN; {
		w, n := getFilterWord(s[i:])
		if n == 0 {
			return nil, errors.New("end of string reached without whitespace delimiter")
		}

		// Sort between stacks of output and operators
		// Determine if operator or operand and sort to stacks
		if fopt, ok := isOperator(w); ok {
			fop, fop_e := new_filter_op_i(fopt)
			if fop_e != nil {
				return nil, fop_e
			}

			// If there are operators in the stack we need to compare priorities
			if ft.operatorStack.Count() > 0 {
				for top_op := ft.operatorStack.Peek(); top_op != nil && fop.p > top_op.p; {
					if err := ft.unwind(); err != nil {
						return nil, err
					}

					top_op = ft.operatorStack.Peek()
				}
			}

			ft.operatorStack.Push(*fop)
		} else {
			// If its not operator then we push it to output stack as a string or number
			if n, err := strconv.Atoi(w); err == nil {
				ft.outputStack.Push(filter_out_i{val: n, valType: NUM})
			} else {
				ft.outputStack.Push(filter_out_i{val: w, valType: STR})
			}
		}

		i += n
	}

	return ft.flush()
}

// Clears the operator stack returning the resulting query
func (ft *FilterTranslator) flush() (*bson.D, error) {
	op_cnt := ft.operatorStack.Count()
	for i := 0; i < op_cnt; i++ {
		if err := ft.unwind(); err != nil {
			return nil, err
		}
	}

	if ft.outputStack.Count() > 1 {
		return nil, errors.New("should only be 1 in output stack from unwinding")
	}

	out := ft.outputStack.Pop()
	return out.val.(*bson.D), nil
}

// Unwinds the top operator and corresponding operand pair. Pushes the output back onto the output stack
func (ft *FilterTranslator) unwind() error {
	if op_cnt, out_cnt := ft.operatorStack.Count(), ft.outputStack.Count(); op_cnt <= 0 || out_cnt < 2 {
		return fmt.Errorf("operator / operand mismatch - have(expect) - operators: %d(>0), operands: %d(>=2)", op_cnt, out_cnt)
	}

	op := ft.operatorStack.Pop()
	y := ft.outputStack.Pop() // Should be property
	x := ft.outputStack.Pop() // Should be operand

	bson, err := fopToBSON(op, x, y)
	if err != nil {
		return err
	}

	ft.outputStack.Push(filter_out_i{val: bson, valType: BSON})
	return nil
}

func fopToBSON(fop *filter_op_i, x *filter_out_i, y *filter_out_i) (*bson.D, error) {
	mop, ok := fop_translation_tb_mongo[fop.t]
	if !ok {
		return nil, fmt.Errorf("%s does not translate to an associated mongodb operator", mop)
	}

	var keyVal string
	var valVal any

	switch fop.t {
	case FOP_OR:
		fallthrough
	case FOP_AND:
		keyVal = string(mop)
		valVal = bson.A{x.val.(*bson.D), y.val.(*bson.D)}
	default:
		keyVal = x.val.(string)
		valVal = bson.D{primitive.E{Key: string(mop), Value: y.val}}
	}

	return &bson.D{primitive.E{Key: keyVal, Value: valVal}}, nil
}

func new_filter_op_i(fopt filter_op_type) (*filter_op_i, error) {
	p, ok := getPriority(fopt)
	if !ok {
		return nil, fmt.Errorf("%s operator does not have a priority", fopt)
	}

	return &filter_op_i{t: fopt, p: p}, nil
}

func getFilterWord(s string) (string, int) {
	sIndx := 0
	eIndx := 0
	inWord := false

	for i, c := range s {
		// Lets skip returning any leading whitespace
		if c == ' ' && !inWord {
			sIndx++
			continue
		} else if c != ' ' && !inWord {
			inWord = true
		} else if c == ' ' && inWord {
			eIndx = i
			break
		}
	}

	// End of string
	if eIndx == 0 {
		eIndx = len(s)
	}

	// End of string case (or no whitespace delimiters)
	retStr := s[sIndx:eIndx]
	return retStr, eIndx
}

func isOperator(s string) (filter_op_type, bool) {
	val, ok := fop_table[s]
	return val, ok
}

func getPriority(fopt filter_op_type) (int, bool) {
	val, ok := fop_priorities[fopt]
	return val, ok
}
