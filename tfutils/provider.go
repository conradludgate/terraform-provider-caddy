package tfutils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Structure represents a type that declares it's structure as a SchemaMap
type Structure interface {
	Schema() SchemaMap
}

// Resource is the type that a terraform resource will implement
type Resource interface {
	CRUD
	Structure
}

// DataSource is the type that a terraform data source will implement
type DataSource interface {
	Data
	Structure
}

// ProviderBuilder is the type that will build into a terraform provider
type ProviderBuilder struct {
	Schema        SchemaMap
	Resources     ResourceMap
	DataSources   DataSourceMap
	ConfigureFunc schema.ConfigureFunc
}

// Build converts the ProviderBuilder into a *schema.Provider
func (pb ProviderBuilder) Build() *schema.Provider {
	return &schema.Provider{
		Schema:         pb.Schema.BuildSchema(),
		ResourcesMap:   pb.Resources.BuildResourcesMap(),
		DataSourcesMap: pb.DataSources.BuildDataSourcesMap(),
		ConfigureFunc:  pb.ConfigureFunc,
	}
}

// ResourceMap represents a map of Resources
type ResourceMap map[string]Resource

// BuildResourcesMap converts a ResourceMap into a map[string]*schema.Resource
func (rm ResourceMap) BuildResourcesMap() map[string]*schema.Resource {
	m := (map[string]Resource)(rm)
	s := make(map[string]*schema.Resource, len(m))
	for k, v := range m {
		s[k] = v.Schema().BuildCRUD(v)
	}
	return s
}

// DataSourceMap represents a map of DataSources
type DataSourceMap map[string]DataSource

// BuildDataSourcesMap converts a DataSourceMap into a map[string]*schema.Resource
func (rm DataSourceMap) BuildDataSourcesMap() map[string]*schema.Resource {
	m := (map[string]DataSource)(rm)
	s := make(map[string]*schema.Resource, len(m))
	for k, v := range m {
		s[k] = v.Schema().BuildDataSource(v)
	}
	return s
}
