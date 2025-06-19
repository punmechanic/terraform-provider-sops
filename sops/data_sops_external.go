package sops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExternalDataSource struct {
	ID     types.String `tfsdk:"id"`
	Source types.String `tfsdk:"source"`
	Format types.String `tfsdk:"input_type"`
	Data   types.Map    `tfsdk:"data"`
	Raw    types.String `tfsdk:"raw"`
}

func (ExternalDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_external"
}

func (e ExternalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if ds := req.Config.Get(ctx, &e); ds.HasError() {
		resp.Diagnostics.Append(ds...)
		return
	}

	source := []byte(e.Source.ValueString())
	format := e.Format.ValueString()
	if err := validateInputType(format); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("input_type"), "bad format", err.Error())
		return
	}

	var fds FileDataSourceModel
	if err := readData(source, format, &fds); err != nil {
		resp.Diagnostics.AddError("failed to read data", err.Error())
		return
	}

	e.ID = types.StringValue(fds.ID)
	e.Raw = types.StringValue(fds.Raw)
	e.Data, _ = types.MapValueFrom(ctx, types.StringType, fds.Data)
	resp.Diagnostics.Append(resp.State.Set(ctx, e)...)
}

func (ExternalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"input_type": schema.StringAttribute{
				Required: true,
			},
			"source": schema.StringAttribute{
				Required: true,
			},
			"data": schema.MapAttribute{
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"raw": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
