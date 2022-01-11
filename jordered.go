package jordered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type OrderedMap struct {
	ordered []element
	iter    int
	raw     interface{}
	rawArr  []interface{}
}

type element struct {
	key   string
	value interface{}
}

// UnmarshalsJSON unmarshals json data and maintains the map order
func (m *OrderedMap) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	m.ordered = []element{}

	// We are only marshalling the first level. Anything deeper we just throw in an interface
	objType := ""
	t, err := dec.Token()
	if err != nil {
		return err
	}

	switch val := t.(type) {
	case json.Delim:
		if strings.ContainsAny(val.String(), "[]") {
			objType = "array"
		} else {
			objType = "object"
		}
	case string, float64, bool, nil, json.Number:
		m.raw = val
		return nil
	}

	if objType == "object" {
		for dec.More() {
			// Read another token as this is the key
			t, err = dec.Token()
			key, ok := t.(string)
			if !ok {
				return fmt.Errorf("Object Key must be a string %T, %v ", t, t)
			}
			omValue := &OrderedMap{}
			err := dec.Decode(&omValue)
			if err != nil {
				return err
			}
			m.ordered = append(m.ordered, element{
				key:   key,
				value: getValue(omValue),
			})
		}
	} else if objType == "array" {
		for dec.More() {
			switch t.(type) {
			case string, float64, bool, nil, json.Number:
				m.rawArr = append(m.rawArr, t)
				t, err = dec.Token()
				if err != nil {
					return err
				}
			default:
				omValue := &OrderedMap{}
				err := dec.Decode(&omValue)
				if err != nil {
					return err
				}
				m.rawArr = append(m.rawArr, getValue(omValue))
			}
		}
	}

	return nil
}

func getValue(omv *OrderedMap) interface{} {
	switch {
	case omv == nil:
		return nil
	case omv.raw != nil:
		return omv.raw
	case omv.rawArr != nil:
		return omv.rawArr
	default:
		return omv
	}
}

func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	if m.raw != nil {
		return json.Marshal(m.raw)
	} else if m.rawArr != nil {
		return json.Marshal(m.rawArr)
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})

	for ii, value := range m.ordered {
		bb, err := json.Marshal(value.value)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buf, `"%s":%s`, value.key, string(bb))
		if ii < len(m.ordered)-1 {
			buf.WriteByte(',')
		}

	}

	buf.Write([]byte{'}'})

	return buf.Bytes(), nil
}

// Next returns true if there is a next item
func (m *OrderedMap) Next() bool {
	return m.iter < len(m.ordered)
}

// Item returns the next item in the map, It returns nil if you have reached the end
func (m *OrderedMap) Item() (string, interface{}) {
	if m.iter >= len(m.ordered) {
		return "", nil
	}
	item := m.ordered[m.iter]
	m.iter++

	return item.key, item.value
}

// Reset restarts the iterator on the map back to the start
func (m *OrderedMap) Reset() {
	m.iter = 0
}

func (m *OrderedMap) Len() int {
	return len(m.ordered)
}

func (m *OrderedMap) Get(key string) (interface{}, bool) {
	for _, item := range m.ordered {
		if item.key == key {
			return item.value, true
		}
	}

	return nil, false
}

func (m *OrderedMap) Set(key string, value interface{}) {
	for ii, item := range m.ordered {
		if item.key == key {
			item.value = value
			m.ordered[ii] = item
			return
		}
	}
	m.Append(key, value)
}

func (m *OrderedMap) Append(key string, value interface{}) {
	if m.ordered == nil {
		m.ordered = []element{}
	}

	m.ordered = append(m.ordered, element{
		key:   key,
		value: value,
	})
}
