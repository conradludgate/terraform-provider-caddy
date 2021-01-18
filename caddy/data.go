package caddy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Data interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
}

type MapData map[string]interface{}

func (rd *MapData) Get(key string) interface{} {
	return (*rd)[key]
}

func (rd *MapData) GetOk(key string) (interface{}, bool) {
	v, ok := (*rd)[key]
	return v, ok
}

func GetString(d Data, key string) string {
	return d.Get(key).(string)
}

func GetStringOk(d Data, key string) *string {
	if v, ok := d.GetOk(key); ok && v.(string) != "" {
		v := v.(string)
		return &v
	}
	return nil
}

func GetBool(d Data, key string) bool {
	return d.Get(key).(bool)
}

func GetInt(d Data, key string) int {
	return d.Get(key).(int)
}

func GetIntOk(d Data, key string) *int {
	if v, ok := d.GetOk(key); ok && v.(int) != 0 {
		v := v.(int)
		return &v
	}
	return nil
}

func GetObject(d Data, key string) MapData {
	return MapData(d.Get(key).(map[string]interface{}))
}

func GetObjectOk(d Data, key string) *MapData {
	if v, ok := d.GetOk(key); ok {
		v := MapData(v.(map[string]interface{}))
		return &v
	}
	return nil
}

func GetObjectSet(d Data, key string) []MapData {
	s := d.Get(key).(*schema.Set).List()
	output := make([]MapData, len(s))
	for i, v := range s {
		output[i] = MapData(v.(map[string]interface{}))
	}
	return output
}

func GetObjectList(d Data, key string) []MapData {
	v := d.Get(key)
	if v == nil {
		return nil
	}
	s := v.([]interface{})
	if len(s) == 0 {
		return nil
	}
	output := make([]MapData, len(s))
	for i, v := range s {
		output[i] = MapData(v.(map[string]interface{}))
	}
	return output
}

func GetStringList(d Data, key string) []string {
	v := d.Get(key)
	if v == nil {
		return nil
	}
	s := v.([]interface{})
	if len(s) == 0 {
		return nil
	}
	output := make([]string, len(s))
	for i, v := range s {
		output[i] = v.(string)
	}
	return output
}

func GetStringMap(d Data, key string) map[string]string {
	s := d.Get(key).(map[string]interface{})
	if len(s) == 0 {
		return nil
	}
	output := make(map[string]string, len(s))
	for i, v := range s {
		output[i] = v.(string)
	}
	return output
}
