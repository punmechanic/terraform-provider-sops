package sops

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTestResourceSopsFile_emptyContentYaml = `
resource "sops_file" "x" {
  content  = ""
  filename = "access-keys.yml"
  pgp {
	fingerprint = "3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A"
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

const configTestResourceSopsFile_withProviderConfig = `
provider "sops" {
  pgp {
	fingerprint = "3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A"
  }
}

resource "sops_file" "x" {
  source      = "%s/test-fixtures/basic-encrypt.yaml"
  filename    = "access-keys.yml"
}`

const configTestResourceSopsFile_withResourceConfig = `
resource "sops_file" "x" {
  source        = "%s/test-fixtures/basic-encrypt.yaml"
  filename      = "access-keys.yml"
  pgp {
	fingerprint = "3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A"
  }
}`

func TestResourceSopsFile_ProviderAndResourceConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(configTestResourceSopsFile_withProviderConfig, wd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.hello", "world"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.integer", "0"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.float", "0.2"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.bool", "true"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.null_value", "null"),
				),
			},
			{
				Config: fmt.Sprintf(configTestResourceSopsFile_withResourceConfig, wd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.hello", "world"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.integer", "0"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.float", "0.2"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.bool", "true"),
					resource.TestCheckResourceAttr("resource.sops_file.x", "data.null_value", "null"),
				),
			},
		},
	})
}
