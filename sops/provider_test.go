package sops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var testAccProviders map[string]func() (tfprotov5.ProviderServer, error)

func init() {
	var p Provider
	testAccProviders = map[string]func() (tfprotov5.ProviderServer, error){
		"sops": providerserver.NewProtocol5WithError(&p),
	}
}

func TestProvider_impl(t *testing.T) {
	var _ provider.Provider = New()
}
