package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/lokkersp/terraform-provider-sops/sops"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), sops.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/registry.terraform.io/lokkersp/sops",
		Debug:   debugMode,
	})

	if err != nil {
		log.Fatal(err)
	}
}
