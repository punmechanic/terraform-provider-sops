package sops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"pgp": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"fingerprint": schema.StringAttribute{
						Description: "The Fingerprint of the PGP key",
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
		Pgp types.Object `tfsdk:"pgp"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &encryptConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var conf encryptConfigModel
	// TODO Only allow one of kms or pgp to be defined
	if !encryptConfig.Kms.IsNull() {
		if ds := unmarshalKmsConf(ctx, encryptConfig.Kms, &conf.Kms); ds.HasError() {
			resp.Diagnostics.Append(ds...)
			return
		}
		conf.EncryptionProvider = "kms"
	}

	if !encryptConfig.Pgp.IsNull() {
		if ds := unmarshalPgpConf(ctx, encryptConfig.Pgp, &conf.Pgp); ds.HasError() {
			resp.Diagnostics.Append(ds...)
			return
		}
		conf.EncryptionProvider = "pgp"
	}

	resp.ResourceData = conf
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
