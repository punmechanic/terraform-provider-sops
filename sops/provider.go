package sops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ provider.Provider = &SopsProvider{}

type SopsProvider struct{}

func New() provider.Provider {
	return &SopsProvider{}
}

func (p *SopsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sops"
}

func (p *SopsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"kms": schema.SingleNestedBlock{
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

func (p *SopsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// TODO: Hacky.
	var encryptConfig struct {
		Kms types.Object `tfsdk:"kms"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &encryptConfig)...)
	if resp.Diagnostics.HasError() || encryptConfig.Kms.IsNull() {
		return
	}

	var kmsConfig kmsConfigSchema
	resp.Diagnostics.Append(encryptConfig.Kms.As(ctx, &kmsConfig, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = KmsConf{
		ARN:     kmsConfig.ARN.ValueString(),
		Profile: kmsConfig.Profile.ValueString(),
	}
}

func (p *SopsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newFileDataSource,
		newExternalDataSource,
	}
}

func (p *SopsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newFileResource,
	}
}
