package hetznerrobot

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAccProviderFactories wires the provider into the acceptance-test harness.
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"hetzner-robot": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

// testAccPreCheck verifies the credentials required for acceptance tests are set.
// Acceptance tests only run when TF_ACC is set and talk to the real Hetzner Robot
// API, so they need a real account.
func testAccPreCheck(t *testing.T) {
	for _, k := range []string{"HETZNERROBOT_USERNAME", "HETZNERROBOT_PASSWORD"} {
		if os.Getenv(k) == "" {
			t.Fatalf("%s must be set for acceptance tests", k)
		}
	}
}

// TestAccServersDataSource exercises the full provider plumbing (configure ->
// live API call -> state) with a read-only servers lookup. It is skipped unless
// TF_ACC is set. Run with:
//
//	TF_ACC=1 HETZNERROBOT_USERNAME=... HETZNERROBOT_PASSWORD=... \
//	  go test ./hetznerrobot/ -run TestAccServersDataSource -v
func TestAccServersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "hetzner-robot" {}

data "hetzner-robot_servers" "all" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hetzner-robot_servers.all", "id"),
				),
			},
		},
	})
}
