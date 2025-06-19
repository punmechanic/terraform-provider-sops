package sops

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FileDataSource struct {
	ID         types.String `tfsdk:"id"`
	InputType  types.String `tfsdk:"input_type"`
	SourceFile types.String `tfsdk:"source_file"`
	Raw        types.String `tfsdk:"raw"`
	Data       types.Map    `tfsdk:"data"`
}

type FileDataSourceModel struct {
	ID         string
	InputType  string
	SourceFile string
	Raw        string
	Data       map[string]string
}

func (fds FileDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_file"
}

func (fds FileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"input_type": schema.StringAttribute{
				Optional: true,
			},
			"source_file": schema.StringAttribute{
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

func (fds FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	diags := req.Config.Get(ctx, &fds)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var m map[string]string
	if diags := fds.Data.ElementsAs(ctx, &m, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	apiModel := FileDataSourceModel{
		SourceFile: fds.SourceFile.ValueString(),
		InputType:  fds.InputType.ValueString(),
		Raw:        fds.Raw.String(),
		Data:       m,
	}

	content, err := os.ReadFile(apiModel.SourceFile)
	if err != nil {
		resp.Diagnostics.AddError("could not read source file", err.Error())
		return
	}

	var format string
	if inputType := apiModel.InputType; inputType != "" {
		format = inputType
	} else {
		switch ext := path.Ext(apiModel.SourceFile); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		case ".env":
			format = "dotenv"
		case ".ini":
			format = "ini"
		default:
			resp.Diagnostics.AddAttributeError(tfpath.Root("source_file"), "unknown format", fmt.Sprintf("Don't know how to decode file with extension %s, set input_type to json, yaml or raw as appropriate", ext))
			return
		}
	}

	if err := validateInputType(format); err != nil {
		resp.Diagnostics.AddAttributeError(tfpath.Root("input_type"), "unknown format", err.Error())
		return
	}

	if err := readData(content, format, &apiModel); err != nil {
		resp.Diagnostics.AddError("read fail", err.Error())
		return
	}

	fds.Raw = types.StringValue(apiModel.Raw)
	fds.Data, diags = types.MapValueFrom(ctx, types.StringType, apiModel.Data)
	fds.SourceFile = types.StringValue(apiModel.SourceFile)
	fds.InputType = types.StringValue(apiModel.InputType)
	fds.ID = types.StringValue(apiModel.ID)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, fds)...)
}
