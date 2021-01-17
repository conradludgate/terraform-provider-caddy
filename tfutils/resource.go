package tfutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SchemaMap is a builder for a *schema.Resource type
type SchemaMap map[string]SchemaBuilder

// BuildSchema converts a SchemaMap into a map[string]*schema.Schema
func (sm SchemaMap) BuildSchema() map[string]*schema.Schema {
	m := (map[string]SchemaBuilder)(sm)
	s := make(map[string]*schema.Schema, len(m))
	for k, v := range m {
		s[k] = v.Build()
	}
	return s
}

// BuildResource converts a SchemaMap into a *schema.Resource
func (sm SchemaMap) BuildResource() *schema.Resource {
	return &schema.Resource{
		Schema: sm.BuildSchema(),
	}
}

// Data represents a read only data source
type Data interface {
	Read(d *schema.ResourceData, m interface{}) error
}

// CRUD represents a crud resource type
type CRUD interface {
	Data
	Create(d *schema.ResourceData, m interface{}) error
	Update(d *schema.ResourceData, m interface{}) error
	Delete(d *schema.ResourceData, m interface{}) error
}

// BuildCRUD creates a new CRUD *schema.Resource type
func (sm SchemaMap) BuildCRUD(crud CRUD) *schema.Resource {
	r := sm.BuildResource()
	r.Create = crud.Create
	r.Delete = crud.Delete
	r.Update = crud.Update
	r.Read = crud.Read
	return r
}

// BuildDataSource creates a new data source *schema.Resource type
func (sm SchemaMap) BuildDataSource(data Data) *schema.Resource {
	r := sm.BuildResource()
	r.Read = data.Read
	return r
}

// IntoSet is the builder version of Set
func (sm SchemaMap) IntoSet() SchemaBuilder {
	return Set(sm)
}

// IntoList is the builder version of List
func (sm SchemaMap) IntoList() SchemaBuilder {
	return List(sm)
}
