package excel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToColumnName(t *testing.T) {
	type testCase struct {
		input  uint
		output string
	}

	testCases := []testCase{
		{1, "A"},
		{2, "B"},
		{26, "Z"},
		{27, "AA"},
		{28, "AB"},
		{52, "AZ"},
		{53, "BA"},
	}

	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := toColumnName(testCase.input)
			assert.Equal(t, testCase.output, actual)
		})
	}

}
