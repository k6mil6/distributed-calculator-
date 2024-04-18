package remaker

import "testing"

type RemakeTest struct {
	subexpression string
	id            int
	result        float64
	expected      string
}

var tests = []RemakeTest{
	{
		subexpression: "{1} * 4234234",
		id:            1,
		result:        40000000,
		expected:      "40000000 * 4234234",
	},
}

func TestRemake(t *testing.T) {
	for _, test := range tests {
		res := Remake(test.subexpression, test.id, test.result)
		if res != test.expected {
			t.Errorf("expected %s, but got: %s", test.expected, res)
		}
	}

}
