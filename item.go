package gospider

import (
	"fmt"
	"encoding/json"
)

/**************************************************************
* struct: Item
**************************************************************/
type Item map[string]interface{}

// Item_Repr implement Data interface
func (item Item) Repr() string {
	return fmt.Sprintf("%+v", item)
}

// Item_Data wraps itself as Data
func (item Item) Data() Data {
	return item
}

// Item_DataList wraps itself as a slice of Data
// useful for those functions who yield only one Item
func (item Item) DataList() []Data {
	return []Data{item}
}

// Item_GetString will access self and assume a string value
// if any error occurs, result is "". so do not use empty string value in Item
func (item Item) GetString(key string) string {
	value, ok := item[key]
	if !ok {
		return ""
	}
	if v, ok := value.(string); !ok {
		return ""
	} else {
		return v
	}
}

// Item_Marshal using json to serialize item
// do not fill value that is not JSON serializable
func (item Item) Marshal() ([]byte, error) {
	return json.Marshal(item)
}

// Item_Unmarshal will deserialize item from []byte
func (item Item) Unmarshal(data []byte) (error) {
	return json.Unmarshal(data, &item)
}
