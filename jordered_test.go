package jordered

import (
	"encoding/json"
	"testing"
)

const (
	mapJson = `{"one":{"hendrik":"sedin"},"two":{"daniel":"sedin"},"three":["vancouver","canucks"],"four":"hockey"}`
	arrJson = `
[
    {
        "hendrik": "sedin" 
    },
    {
        "daniel": "sedin"
    },
    [
        "vancouver",
        "canucks"
    ],
    "hockey"
]
`
)

func TestOrderedMap(t *testing.T) {
	om := &OrderedMap{}

	err := json.Unmarshal([]byte(mapJson), om)
	if err != nil {
		t.Error(err)
	}

	bb, err := json.Marshal(om)
	if err != nil {
		t.Error(err)
	}

	if string(bb) != mapJson {
		t.Errorf("Bad json marshall expected %s, got %s", mapJson, string(bb))
	}

	keys := []string{"one", "two", "three", "four"}
	values := [][]byte{
		[]byte(`{"hendrik":"sedin"}`),
		[]byte(`{"daniel":"sedin"}`),
		[]byte(`["vancouver","canucks"]`),
		[]byte(`"hockey"`),
	}

	// Test iterator and order
	ii := 0
	for om.Next() {
		key, value := om.Item()
		if keys[ii] != key {
			t.Errorf("Item does not have the correct key expected: %s got: %s", keys[ii], key)
		}

		valB, err := json.Marshal(value)
		if err != nil {
			t.Error(err)
		}
		if string(values[ii]) != string(valB) {
			t.Errorf("Item does not have the correct values expected: %s got: %s", string(values[ii]), string(valB))
		}

		ii++
	}
	if ii != 4 {
		t.Errorf("Did not find all the values (or two many). Expected %v, Got %v", 4, ii)
	}

	// Test Reset
	om.Reset()
	ii = 0
	for om.Next() {
		key, value := om.Item()
		if keys[ii] != key {
			t.Errorf("Item does not have the correct key expected: %s got: %s", keys[ii], key)
		}

		valB, err := json.Marshal(value)
		if err != nil {
			t.Error(err)
		}
		if string(values[ii]) != string(valB) {
			t.Errorf("Item does not have the correct values expected: %s got: %s", string(values[ii]), string(valB))
		}

		ii++
	}
	if ii != 4 {
		t.Errorf("Did not find all the values (or two many). Expected %v, Got %v", 4, ii)
	}
	om.Reset()
	// Test append
	om.Append("five", 127)
	keys = append(keys, "five")
	values = append(values, []byte(`127`))
	ii = 0
	for om.Next() {
		key, value := om.Item()
		if keys[ii] != key {
			t.Errorf("Item does not have the correct key expected: %s got: %s", keys[ii], key)
		}

		valB, err := json.Marshal(value)
		if err != nil {
			t.Error(err)
		}
		if string(values[ii]) != string(valB) {
			t.Errorf("Item does not have the correct values expected: %s got: %s", string(values[ii]), string(valB))
		}

		ii++
	}
	if ii != 5 {
		t.Errorf("Did not find all the values (or two many). Expected %v, Got %v", 5, ii)
	}

	// Test get
	value, ok := om.Get("two")
	if !ok {
		t.Errorf("Key two does not exist, it should")
	}
	valB, err := json.Marshal(value)
	if err != nil {
		t.Error(err)
	}
	if string(values[1]) != string(valB) {
		t.Errorf("Item does not have the correct values expected: %s got: %s", string(values[1]), string(valB))
	}

	// Test set
	om.Set("two", "foo")
	value, ok = om.Get("two")
	if !ok {
		t.Errorf("Key two does not exist, it should")
	}
	if "foo" != value.(string) {
		t.Errorf("Item does not have the correct values expected: foo got: %s", value)
	}

	om = &OrderedMap{}

	err = json.Unmarshal([]byte(arrJson), om)
	if err == nil {
		t.Error("Marshalling should have failed. It succeeded")
	}

}
