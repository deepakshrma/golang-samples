package number

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type TestData [][]float64
type Case struct {
	name     string
	a        Num
	b        Num
	expected Num
}

func init() {
	/* load test data */
	fmt.Println("Init")
}
func prepareData() []Case {
	jsonFile, err := os.Open("./../testdata/data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var testdata TestData
	var cases []Case
	json.Unmarshal(byteValue, &testdata)
	for _, data := range testdata {
		cases = append(cases, Case{
			"add+++",
			Num(data[0]),
			Num(data[1]),
			Num(data[2]),
		})
	}
	return cases
}
func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	t.Log("setup sub test")
	return func(t *testing.T) {
		t.Log("teardown sub test")
	}
}
func TestAdd(t *testing.T) {
	calc := Num(1.0)

	if calc.Add(3) != 4 {
		t.Errorf("%f + %f != %f", 1.0, 3.0, calc.Add(3))
	}
}
func TestAddDataTable(t *testing.T) {
	testdata := prepareData()
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for _, tc := range testdata {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := setupSubTest(t)
			defer teardownSubTest(t)
			result := tc.a.Add(tc.b)
			if result != tc.expected {
				t.Fatalf("expected sum %v, but got %v", tc.expected, result)
			}
		})
	}
}
