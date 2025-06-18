package sops

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var sourceFile = resourceSourceFile()

// Fixes a bug where sops_file would crash if filename was set to .yml and content was empty string.
func TestResourceSopsFile_ReturnsCouldNotReadInputFileIfYmlFileIsEmpty(t *testing.T) {
	var instanceState terraform.InstanceState
	data := sourceFile.Data(&instanceState)
	data.Set("content", "")
	data.Set("filename", "access-keys.yml")

	ds := sourceFile.CreateContext(context.Background(), data, &EncryptConfig{})
	if !ds.HasError() {
		t.Error("No errors found, but expected exactly one")
	}

	first := ds[0]
	expected := "provided content was empty\n"
	if first.Summary != expected {
		t.Errorf("Error %d Summary was %q instead of %q", 0, first.Summary, expected)
	}
}
