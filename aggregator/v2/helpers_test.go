package aggregator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumToChars(t *testing.T) {
	type testCase struct {
		input  uint
		output string
	}

	testCases := []testCase{
		{1, "a"},
		{2, "b"},
		{26, "z"},
		{27, "aa"},
		{28, "ab"},
		{52, "az"},
		{53, "ba"},
	}

	for i, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := numToChars(testCase.input)
			assert.Equal(t, testCase.output, actual)
		})
	}

}
