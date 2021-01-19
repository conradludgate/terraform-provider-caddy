package tfutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SchemaBuilder is a utility type that makes it easy to create new schema objects
type SchemaBuilder struct {
	s schema.Schema
}

// Required will set schema as required
func (b SchemaBuilder) Required() SchemaBuilder {
	b.s.Required = true
	b.s.Optional = false

	return b
}

// Optional will set schema as optional
func (b SchemaBuilder) Optional() SchemaBuilder {
	b.s.Required = false
	b.s.Optional = true

	return b
}

// Computed will set schema as computed
func (b SchemaBuilder) Computed() SchemaBuilder {
	b.s.Computed = true
	return b
}

// Build converts the SchemaBuiler into a Schema
func (b SchemaBuilder) Build() *schema.Schema {
	return &b.s
}

// String creates a new string schema
func String() SchemaBuilder {
	return NewSchema().WithType(schema.TypeString)
}

// Int creates a new string schema
func Int() SchemaBuilder {
	return NewSchema().WithType(schema.TypeInt)
}

// Bool creates a new string schema
func Bool() SchemaBuilder {
	return NewSchema().WithType(schema.TypeBool)
}

// NewSchema creates a new empty schema builder
func NewSchema() SchemaBuilder {
	return SchemaBuilder{}
}

// WithType sets the type of the schema
func (b SchemaBuilder) WithType(typ schema.ValueType) SchemaBuilder {
	b.s.Type = typ
	return b
}

// WithElem sets the elem of the schema
func (b SchemaBuilder) WithElem(elem interface{}) SchemaBuilder {
	b.s.Elem = elem
	return b
}

// Default sets the default element of the schema and sets the schema to optional
func (b SchemaBuilder) Default(v interface{}) SchemaBuilder {
	b.s.Default = v
	b.s.Optional = true
	return b
}

// List converts the type into a list. Must be called after String, Int etc
func (b SchemaBuilder) List() SchemaBuilder {
	return b.WithElem(&schema.Schema{Type: b.s.Type}).WithType(schema.TypeList)
}

// Map converts the type into a map. Must be called after String, Int etc
func (b SchemaBuilder) Map() SchemaBuilder {
	return b.WithElem(&schema.Schema{Type: b.s.Type, Elem: b.s.Elem}).WithType(schema.TypeMap)
}

// Set converts the type into a set. Must be called after String, Int etc
func (b SchemaBuilder) Set() SchemaBuilder {
	return b.WithElem(&schema.Schema{Type: b.s.Type}).WithType(schema.TypeSet)
}

// ListOf creates a new schema list over the resource
func ListOf(r Structure) SchemaBuilder {
	return NewSchema().WithType(schema.TypeList).WithElem(r.Schema().BuildResource())
}

// SetOf creates a new schema set over the resource
func SetOf(r Structure) SchemaBuilder {
	return NewSchema().WithType(schema.TypeSet).WithElem(r.Schema().BuildResource())
}

// Set creates a schema set type over the map of sub schemas
func Set(s SchemaMap) SchemaBuilder {
	return NewSchema().WithType(schema.TypeSet).WithElem(s.BuildResource())
}

// List creates a schema list type over the map of sub schemas
func List(s SchemaMap) SchemaBuilder {
	return NewSchema().WithType(schema.TypeList).WithElem(s.BuildResource())
}

// ConflictsWith adds the given keys to the schema's conflicts with array
func (b SchemaBuilder) ConflictsWith(key ...string) SchemaBuilder {
	b.s.ConflictsWith = append(b.s.ConflictsWith, key...)
	return b
}

// MaxItems sets the maximum number of items allowed in this set/list
func (b SchemaBuilder) MaxItems(n int) SchemaBuilder {
	b.s.MaxItems = n
	return b
}

func (b SchemaBuilder) SetFunc(f schema.SchemaSetFunc) SchemaBuilder {
	b.s.Set = f
	return b
}
