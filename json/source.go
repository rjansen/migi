package environment

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/rjansen/migi"
)

type source struct {
	reader  io.Reader
	options map[string]interface{}
}

func (e *source) Load() error {
	return json.NewDecoder(e.reader).Decode(&e.options)
}

func (e *source) lookup(name string) (interface{}, error) {
	value, ok := e.options[name]
	if !ok {
		return nil, migi.NewOptionNotFound(name)
	}

	return value, nil
}

func (e *source) String(name string) (string, error) {
	value, err := e.lookup(name)
	if err != nil {
		return "", err
	}

	strValue, is := value.(string)
	if !is {
		return "", migi.NewOptionInvalidType(name, value, "string")
	}

	return strValue, nil
}

func (e *source) Int(name string) (int, error) {
	value, err := e.lookup(name)
	if err != nil {
		return 0, err
	}

	switch rawValue := value.(type) {
	case string:
		intValue, err := strconv.ParseInt(rawValue, 10, 32)
		if err != nil {
			return 0, err
		}
		return int(intValue), nil
	case float64:
		return int(rawValue), nil
	default:
		return 0, migi.NewOptionInvalidType(name, value, "int")
	}
}

func (e *source) Float(name string) (float32, error) {
	value, err := e.lookup(name)
	if err != nil {
		return 0, err
	}

	switch rawValue := value.(type) {
	case string:
		floatValue, err := strconv.ParseFloat(rawValue, 32)
		if err != nil {
			return 0, err
		}
		return float32(floatValue), nil
	case float64:
		return float32(rawValue), nil
	default:
		return 0, migi.NewOptionInvalidType(name, value, "float")
	}
}

func (e *source) Bool(name string) (bool, error) {
	value, err := e.lookup(name)
	if err != nil {
		return false, err
	}

	switch rawValue := value.(type) {
	case string:
		boolValue, err := strconv.ParseBool(rawValue)
		if err != nil {
			return false, err
		}
		return boolValue, nil
	case bool:
		return rawValue, nil
	default:
		return false, migi.NewOptionInvalidType(name, value, "float")
	}
}

func (e *source) Time(name string) (time.Time, error) {
	value, err := e.lookup(name)
	if err != nil {
		return time.Time{}, err
	}

	switch rawValue := value.(type) {
	case string:
		timeValue, err := time.Parse(time.RFC3339, rawValue)
		if err != nil {
			return time.Time{}, err
		}
		return timeValue, nil
	case time.Time:
		return rawValue, nil
	default:
		return time.Time{}, migi.NewOptionInvalidType(name, value, "time.Time")
	}
}

func (e *source) Duration(name string) (time.Duration, error) {
	value, err := e.lookup(name)
	if err != nil {
		return time.Duration(0), err
	}

	switch rawValue := value.(type) {
	case string:
		durationValue, err := time.ParseDuration(rawValue)
		if err != nil {
			return time.Duration(0), err
		}
		return durationValue, nil
	case time.Duration:
		return rawValue, nil
	default:
		return time.Duration(0), migi.NewOptionInvalidType(name, value, "time.Duration")
	}
}

func NewSource(reader io.Reader) *source {
	return &source{
		reader:  reader,
		options: make(map[string]interface{}),
	}
}
