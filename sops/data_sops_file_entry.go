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

type FileEntryDataSource struct {
	ID         types.String `tfsdk:"source_file"`
	SourceFile types.String `tfsdk:"source_file"`
	AgeKeyFile types.String `tfsdk:"age_key_file"`
	InputType  types.String `tfsdk:"input_type"`
	DataKey    types.String `tfsdk:"data_key"`
	Data       types.String `tfsdk:"data"`
	Raw        types.String `tfsdk:"raw"`
	Yaml       types.String `tfsdk:"yaml"`
	Map        types.Map    `tfsdk:"map"`
}

func (FileEntryDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, req *datasource.MetadataResponse) {
	req.TypeName = "sops_file_entry"
}

func (f FileEntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	if ds := req.Config.Get(ctx, &f); ds.HasError() {
		res.Diagnostics.Append(ds...)
		return
	}

	sourceFile := f.SourceFile.ValueString()
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		res.Diagnostics.AddError("cannot read file", fmt.Sprintf("cannot read file %s: %s", sourceFile, err))
		return
	}

	if !f.AgeKeyFile.IsNull() {
		ageKeyFile := f.AgeKeyFile.ValueString()
		if ageKeyFile != "" {
			envVarName := "SOPS_AGE_KEY_FILE"
			err := os.Setenv(envVarName, ageKeyFile)
			if err != nil {
				log.Errorf("fail to set environment variable %s.Error is %s", envVarName, err)
			}
		}
	}

	var format string
	if inputType := f.InputType.ValueString(); inputType != "" {
		format = inputType
	} else {
		switch ext := path.Ext(sourceFile); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		case ".env":
			format = "dotenv"
		case ".ini":
			format = "ini"
		default:
			res.Diagnostics.AddAttributeError(tfpath.Root("source_file"), "unknown file extension", fmt.Sprintf("Don't know how to decode file with extension %s, set input_type to json, yaml or raw as appropriate", ext))
			return
		}
	}

	if err := validateInputType(format); err != nil {
		res.Diagnostics.AddAttributeError(tfpath.Root("input_type"), "bad input type", err.Error())
		return
	}

	var dk readDataKeyModel
	if err := readDataKey(content, format, f.DataKey.ValueString(), &dk); err != nil {
		res.Diagnostics.AddError("could not read data key", err.Error())
		return
	}

	f.ID = types.StringValue(dk.ID)
	f.Raw = types.StringValue(dk.Raw)
	f.Data = types.StringValue(dk.Data)
	f.Yaml = types.StringValue(dk.Yaml)
	f.Map, _ = types.MapValueFrom(ctx, types.StringType, dk.Data)
	res.Diagnostics.Append(res.State.Set(ctx, f)...)
}

func (FileEntryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"data_key": schema.StringAttribute{
				Required: true,
			},
			"age_key_file": schema.StringAttribute{
				Optional: true,
			},
			"data": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"yaml": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"map": schema.MapAttribute{
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
