package caddy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// slice interface to string
func sitos(s []interface{}) []string {
	output := make([]string, len(s))
	for i, v := range s {
		output[i] = v.(string)
	}
	return output
}

// slice interface to string
func sstoi(s []string) []interface{} {
	output := make([]interface{}, len(s))
	for i, v := range s {
		output[i] = v
	}
	return output
}

// GetOkS returns the string value at key if set, otherwise nil
func GetOkS(r *schema.ResourceData, key string) *string {
	if v, ok := r.GetOk(key); ok {
		s := v.(string)
		return &s
	}
	return nil
}

// GetOkB returns the bool value at key if set, otherwise nil
func GetOkB(r *schema.ResourceData, key string) *bool {
	if v, ok := r.GetOk(key); ok {
		b := v.(bool)
		return &b
	}
	return nil
}
