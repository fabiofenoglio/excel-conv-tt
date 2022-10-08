package database

import "strings"

var (
	knownOperatorMap map[string]KnownOperator
)

func init() {
	knownOperatorMap = make(map[string]KnownOperator)

	registerKnownOperators(KnownOperator{
		Code:            "a",
		Name:            "Op. A",
		BackgroundColor: "#E0EBF5",
	})
}

func registerKnownOperators(o ...KnownOperator) {
	for _, obj := range o {
		if obj.Code == "" {
			panic("known operator must have a code")
		}
		knownOperatorMap[strings.ToLower(obj.Code)] = obj
	}

}

func GetKnownOperator(code string) (KnownOperator, bool) {
	res, ok := knownOperatorMap[code]
	return res, ok
}

func GetKnownOperators() []KnownOperator {
	out := make([]KnownOperator, 0, len(knownOperatorMap))
	for _, v := range knownOperatorMap {
		out = append(out, v)
	}
	return out
}
