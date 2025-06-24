package sops

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTestResourceSopsFile_emptyContentYaml = `
resource "sops_file" "x" {
  content  = ""
  filename = "access-keys.yml"
  kms {
	arn = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
  }
}`

// Fixes a bug where sops_file would crash if filename was set to .yml and content was empty string.
func TestResourceSopsFile_ReturnsCouldNotReadInputFileIfYmlFileIsEmpty(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configTestResourceSopsFile_emptyContentYaml,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("resource.sops_file.x", "content", ""),
				),
				ExpectError: regexp.MustCompile("provided content was empty"),
			},
		},
	})
}
