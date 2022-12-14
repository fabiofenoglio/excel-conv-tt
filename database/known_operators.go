package database

import "strings"

var (
	knownOperatorMap      map[string]KnownOperator
	knownOperatorAliasMap map[string]string
)

func init() {
	knownOperatorMap = make(map[string]KnownOperator)
	knownOperatorAliasMap = make(map[string]string)

	registerKnownOperators(
		KnownOperator{
			Code:            "emanuele",
			Name:            "Emanuele",
			BackgroundColor: "#A0FC4E",
			Aliases:         []string{"ema", "emanuelebalboni", "balboni"},
		},
		KnownOperator{
			Code:            "jonida",
			Name:            "Jo",
			BackgroundColor: "#FDD1C0",
			Aliases:         []string{"jo", "jonidahalo", "jay", "halo"},
		},
		KnownOperator{
			Code:            "marco",
			Name:            "Marco",
			BackgroundColor: "#2B66B3",
			Aliases:         []string{"marcobrusaferro", "brusaferro", "brusa", "brusamarco"},
		},
		KnownOperator{
			Code:            "roberta",
			Name:            "Roberta",
			BackgroundColor: "#75FBFC",
			Aliases:         []string{"robi", "boccomino"},
		},
		KnownOperator{
			Code:            "lorenzo",
			Name:            "Lorenzo",
			BackgroundColor: "#E7C656",
			Aliases:         []string{"lorenzocolombo", "colombo"},
		},
		KnownOperator{
			Code:            "eleonora",
			Name:            "Eleonora",
			BackgroundColor: "#B9342A",
			Aliases:         []string{"ele"},
		},
		KnownOperator{
			Code:            "simonaromaniello",
			Name:            "Simona Ro.",
			BackgroundColor: "#B2B2B2",
			Aliases:         []string{"ro", "simoro", "romaniello"},
		},
		KnownOperator{
			Code:            "simonarachetto",
			Name:            "Simona Ra.",
			BackgroundColor: "#964B7C",
			Aliases:         []string{"ra", "simora", "rachetto"},
		},
	)
}

func registerKnownOperators(o ...KnownOperator) {
	for _, obj := range o {
		if obj.Code == "" {
			panic("known operator must have a code")
		}
		knownOperatorMap[strings.ToLower(obj.Code)] = obj

		for _, alias := range obj.Aliases {
			knownOperatorAliasMap[alias] = obj.Code
		}
	}
}

func GetKnownOperator(code string) (KnownOperator, bool) {
	if aliasOf, isAlias := knownOperatorAliasMap[code]; isAlias {
		return GetKnownOperator(aliasOf)
	}

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
