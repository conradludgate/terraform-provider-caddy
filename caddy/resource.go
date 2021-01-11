package caddy

// import (
// 	"log"
// 	"strconv"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func embed(embeded bool, key string, r *schema.Resource) {
// 	if !embeded {
// 		r.Schema[key] = &schema.Schema{
// 			Type:     schema.TypeString,
// 			Required: true,
// 		}
// 	} else {
// 		r.Schema["id"] = &schema.Schema{
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		}
// 	}
// }

// // ResourceData is similar to schema.ResourceData but has support for prefixes
// type ResourceData struct {
// 	d      *schema.ResourceData
// 	prefix string
// }

// // RangeSet ranges over the schema.Set value in `key` as a new ResourceData item
// func RangeSet(d *schema.ResourceData, key string, fn func(*ResourceData) error) error {
// 	return WithPrefix(d, "").RangeSet(key, fn)
// }

// // RangeSet ranges over the schema.Set value in `key` as a new ResourceData item
// func (r *ResourceData) RangeSet(key string, fn func(*ResourceData) error) error {
// 	s := r.Get(key)
// 	if s == nil {
// 		return nil
// 	}
// 	set := s.(*schema.Set)
// 	for _, elem := range set.List() {
// 		hash := strconv.Itoa(set.F(elem))
// 		if err := fn(r.WithPrefix(key + "." + hash + ".")); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // WithPrefix creates a new ResourceData with the new prefix
// func WithPrefix(r *schema.ResourceData, prefix string) *ResourceData {
// 	return &ResourceData{r, prefix}
// }

// // WithPrefix creates a new ResourceData with the new prefix
// func (r *ResourceData) WithPrefix(prefix string) *ResourceData {
// 	return &ResourceData{
// 		r.d,
// 		r.prefix + prefix,
// 	}
// }

// // Get returns the value at key
// func (r *ResourceData) Get(key string) interface{} {
// 	return r.d.Get(r.prefix + key)
// }

// // GetS returns the string value at key
// func (r *ResourceData) GetS(key string) string {
// 	return r.d.Get(r.prefix + key).(string)
// }

// // GetOk gets the value at key, returning OK if the value was found
// func (r *ResourceData) GetOk(key string) (interface{}, bool) {
// 	return r.d.GetOk(r.prefix + key)
// }

// // GetOkS returns the string value at key if set, otherwise nil
// func (r *ResourceData) GetOkS(key string) *string {
// 	if v, ok := r.GetOk(key); ok {
// 		s := v.(string)
// 		return &s
// 	}
// 	return nil
// }

// // GetOkB returns the bool value at key if set, otherwise nil
// func (r *ResourceData) GetOkB(key string) *bool {
// 	if v, ok := r.GetOk(key); ok {
// 		b := v.(bool)
// 		return &b
// 	}
// 	return nil
// }

// // Set sets the value at key
// func (r *ResourceData) Set(key string, value interface{}) error {
// 	return r.d.Set(r.prefix+key, value)
// }

// // HasChange detects if key has changed
// func (r *ResourceData) HasChange(key string) bool {
// 	return r.d.HasChange(r.prefix + key)
// }

// // GetChange gets the changes of key
// func (r *ResourceData) GetChange(key string) (interface{}, interface{}) {
// 	return r.d.GetChange(r.prefix + key)
// }

// // Id returns the resource ID
// func (r *ResourceData) Id() string {
// 	if r.prefix == "" {
// 		return r.d.Id()
// 	}
// 	return r.Get("id").(string)
// }

// // SetId sets the resource ID
// func (r *ResourceData) SetId(id string) {
// 	log.Println("SetID", r, id)
// 	if r.prefix == "" {
// 		r.d.SetId(id)
// 	} else {
// 		r.Set("id", id)
// 	}
// }
