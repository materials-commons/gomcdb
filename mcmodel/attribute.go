package mcmodel

import "encoding/json"

const (
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
	ID    int
	Name  string
	Val   string
	Value AttributeValue
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
	var val map[string]interface{}
	if err := json.Unmarshal([]byte(a.Val), &val); err != nil {
		return err
	}
	switch val["val"].(type) {
	case int:
		a.Value.ValueType = ValueTypeInt
		a.Value.ValueInt = val["val"].(int64)
	case []int:
		a.Value.ValueType = ValueTypeArrayOfInt
		a.Value.ValueArrayOfInt = val["val"].([]int64)
	case float32, float64:
		a.Value.ValueType = ValueTypeFloat
		a.Value.ValueFloat = val["val"].(float64)
	case []float32, []float64:
		a.Value.ValueType = ValueTypeArrayOfFloat
		a.Value.ValueArrayOfFloat = val["val"].([]float64)
	case string:
		return nil
	case []string:
		return nil
	case map[interface{}]interface{}:
		return nil
	case []map[interface{}]interface{}:
	default:
	}

	return nil
}

func (a Attribute) GetValue() interface{} {
	return nil
}
