package jordered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type OrderedMap struct {
	ordered []element
	iter    int
}

type element struct {
	key   string
	value interface{}
}

// UnmarshalsJSON unmarshals json data and maintains the map order
func (m *OrderedMap) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	m.ordered = []element{}

	// We are only marshalling the first level. Anything deeper we just throw in an interface
	depth := 0

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if depth == 0 {
			switch val := t.(type) {
			case json.Delim:
				if strings.ContainsAny(val.String(), "[]") {
					// This is an array not a map.
					return fmt.Errorf("Unable to unmarshall array into map")
				}
				depth++
				continue
			case string, float64, bool, nil, json.Number:
				return fmt.Errorf("Unable to unmarshall value into map")
			}
		}

		if depth == 1 {
			key, ok := t.(string)
			if !ok {
				end, ok := t.(json.Delim)
				if !ok {
					return fmt.Errorf("Object Key must be a string %T, %v ", t, t)
				}
				if end.String() == "}" {
					// We are done
					break
				}
			}
			var value interface{}
			err := dec.Decode(&value)
			if err != nil {
				return err
			}
			m.ordered = append(m.ordered, element{
				key:   key,
				value: value,
			})
		}
	}

	return nil
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
