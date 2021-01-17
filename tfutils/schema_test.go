package tfutils_test

import (
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	s := tfutils.String().List().Build()
	assert.Equal(t, &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
	}, s)
}
