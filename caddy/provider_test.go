package caddy

import (
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddy/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// var testAccProvider *schema.Provider
// var testAccProviders map[string]*schema.Provider

// func TestMain(t *testing.M) {
// 	if os.Getenv("TF_ACC") != "" {
// 		testAccProvider = Provider()
// 		testAccProviders = map[string]*schema.Provider{
// 			"caddy": testAccProvider,
// 		}
// 	}

// 	t.Run()
// }

func UnitTest(t *testing.T, setup func(*mocks.Client), steps ...resource.TestStep) {
	t.Helper()

	testAccProvider := Provider()
	caddyMock := &mocks.Client{}
	setup(caddyMock)
	testAccProvider.ConfigureFunc = func(_ *schema.ResourceData) (interface{}, error) {
		return caddyMock, nil
	}
	testAccProviders := map[string]*schema.Provider{
		"caddy": testAccProvider,
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps:      steps,
	})

	caddyMock.AssertExpectations(t)
}
