package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScaffolding(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSealedSecretsSecret,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sealedsecrets_secret.foo", "value", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceSealedSecretsSecret = `
resource "sealedsecrets_secret" "foo" {
  encrypted_value = "bar"
}
`
