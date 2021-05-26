package mcmodel

import (
	"encoding/json"
)

const (
	ValueTypeUnset          = 0
	ValueTypeInt            = 1
	ValueTypeFloat          = 2
	ValueTypeString         = 3
	ValueTypeComplex        = 4
	ValueTypeArrayOfInt     = 5
	ValueTypeArrayOfFloat   = 6
	ValueTypeArrayOfString  = 7
	ValueTypeArrayOfComplex = 8
)

type Attribute struct {
	ID               int
	Name             string
	AttributableID   int
	AttributableType string
	Val              string
	Value            AttributeValue `gorm:"-"`
}

type AttributeValue struct {
	ValueType           int
	ValueInt            int64
	ValueFloat          float64
	ValueComplex        map[string]interface{}
	ValueString         string
	ValueArrayOfInt     []int64
	ValueArrayOfFloat   []float64
	ValueArrayOfString  []string
	ValueArrayOfComplex []map[string]interface{}
}

func (a *Attribute) LoadValue() error {
	if a.Value.ValueType != ValueTypeUnset {
		return nil
	}

	var val map[string]interface{}
	if err := json.Unmarshal([]byte(a.Val), &val); err != nil {
		return err
	}
	//fmt.Printf("%+v\n", val)
	switch val["value"].(type) {
	case int:
		a.Value.ValueType = ValueTypeInt
		a.Value.ValueInt = val["value"].(int64)
	case []int:
		a.Value.ValueType = ValueTypeArrayOfInt
		a.Value.ValueArrayOfInt = val["value"].([]int64)
	case float32, float64:
		a.Value.ValueType = ValueTypeFloat
		a.Value.ValueFloat = val["value"].(float64)
	case []float32, []float64:
		a.Value.ValueType = ValueTypeArrayOfFloat
		a.Value.ValueArrayOfFloat = val["value"].([]float64)
	case string:
		a.Value.ValueType = ValueTypeString
		a.Value.ValueString = val["value"].(string)
	case []string:
		// support later
	case map[interface{}]interface{}:
		// support later
	case []map[interface{}]interface{}:
		// support later
	default:
		//fmt.Printf("Unknown cast type for attribute %s\n", a.Name)
		// What to do here?
	}

	return nil
}

func (a Attribute) GetValue() interface{} {
	return nil
}
