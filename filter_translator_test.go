package quang

import (
	"bytes"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var filter_test_data = [][2]string{
	{
		`Title LE /Test/ AND Total EQ 1`,
		`{"$and":[{"Title":{"$regex":"/Test/"}},{"Total":{"$eq":1}}]}`,
	},
}

func TestTranslateToMongoDB(t *testing.T) {
	for _, tst := range filter_test_data {
		ft := NewFilterTranslator()
		filterStr := tst[0]
		expectedBsonStr := tst[1]

		receivedBSON, err := ft.TranslateToMongo(filterStr)
		if err != nil {
			FailWithError(t, "", err)
		}

		comparePrintBSON(t, receivedBSON, expectedBsonStr)
	}
}

func comparePrintBSON(t *testing.T, xBson *bson.D, y string) {
	if comp, comp_e := compareBSON(xBson, y); !comp {
		jsonReceived, _ := bson.MarshalExtJSON(*xBson, false, false)

		FailWithMessagef(t, "BSON\nexpected:\n%s\nreceived:\n%s\nwith error:\n%s", y, jsonReceived, comp_e.Error())
	}
}

// Returns true if the same
func compareBSON(x *bson.D, y string) (bool, error) {
	var xB, yB []byte
	var err error

	xB, err = bson.Marshal(x)
	if err == nil {
		yBson := &bson.D{}
		if err = bson.UnmarshalExtJSON([]byte(y), false, yBson); err != nil {
			return false, err
		}

		yB, err = bson.Marshal(yBson)
		return bytes.Equal(xB, yB), err
	}

	return false, err
}

func FailWithError(t *testing.T, description string, err error) {
	t.Fatalf("%s\n%s\n", description, err.Error())
}

func FailWithMessage(t *testing.T, message string) {
	t.Fatalf("%s\n", message)
}

func FailWithMessagef(t *testing.T, format string, v ...any) {
	t.Fatalf(format, v...)
}
