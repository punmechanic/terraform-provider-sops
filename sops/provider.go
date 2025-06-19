package sops

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct {
	KMS types.Object `tfsdk:"kms"`
}

func New() provider.Provider {
	return &Provider{}
}

func (Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sops"
}

func (Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"kms": schema.SingleNestedBlock{
				Description: "The default KMS configuration for all resources configured with this provider",
				Attributes: map[string]schema.Attribute{
					"arn": schema.StringAttribute{
						Description: "The ARN of the KMS key",
						Optional:    true,
					},
					"profile": schema.StringAttribute{
						Description: "The AWS Profile to use when retrieving the key",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (p Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if diags := req.Config.Get(ctx, &p); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var encConf EncryptConfig
	if diags := unmarshalKmsConf(ctx, p.KMS, &encConf.Kms); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		fmt.Println("failed to init kms")
	}

	resp.ResourceData = encConf
}

func (Provider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &FileDataSource{} },
		func() datasource.DataSource { return &FileEntryDataSource{} },
		func() datasource.DataSource { return &ExternalDataSource{} },
	}
}

func (Provider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &FileResource{} },
	}
}
