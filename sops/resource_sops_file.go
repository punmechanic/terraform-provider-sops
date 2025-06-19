package sops

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/getsops/sops/v3/aes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FileResource struct {
	// ProviderConfig is the default configuration from the Provider for encryption.
	// This may be a zero value.
	ProviderConfig EncryptConfig
}

type FileResourceIdentity struct {
	ID types.String `tfsdk:"id"`
}

type FileResourceModel struct {
	SensitiveContent    types.String `tfsdk:"sensitive_content"`
	ContentBase64       types.String `tfsdk:"content_base64"`
	Source              types.String `tfsdk:"source"`
	Content             types.String `tfsdk:"content"`
	FilePermission      types.String `tfsdk:"file_permission"`
	DirectoryPermission types.String `tfsdk:"directory_permission"`
	Filename            types.String `tfsdk:"filename"`
	EncryptedRegex      types.String `tfsdk:"encrypted_regex"`
	Kms                 types.Object `tfsdk:"kms"`
}

type FileResourceAPIModel struct {
	Filename            string
	SensitiveContent    string
	ContentBase64       string
	Source              string
	Content             string
	FilePermission      string
	DirectoryPermission string
	EncryptedRegex      string
	Kms                 *KmsConf
}

func (f *FileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if cfg, ok := req.ProviderData.(EncryptConfig); ok {
		f.ProviderConfig = cfg
	}
}

func (FileResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var filename string
	if ds := req.State.GetAttribute(ctx, tfpath.Root("filename"), &filename); ds.HasError() {
		res.Diagnostics.Append(ds...)
		return
	}

	os.Remove(filename)
}

func (FileResource) Metadata(_ context.Context, _ resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = "sops_file"
}

func (FileResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	// If the output file doesn't exist, mark the resource for creation.
	var filename string
	var id FileResourceIdentity
	if ds := req.State.GetAttribute(ctx, tfpath.Root("filename"), &filename); ds.HasError() {
		res.Diagnostics.Append(ds...)
		return
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Unset ID
		res.Diagnostics.Append(res.Identity.Set(ctx, &FileResourceIdentity{ID: types.StringValue("")})...)
		return
	}

	// Verify that the content of the destination file matches the content we
	// expect. Otherwise, the file might have been modified externally and we
	// must reconcile.
	outputContent, err := os.ReadFile(filename)
	if err != nil {
		res.Diagnostics.AddError(fmt.Sprintf("failed to read %s", filename), err.Error())
		return
	}

	outputChecksum := sha1.Sum(outputContent)
	expectedID := hex.EncodeToString(outputChecksum[:])
	if ds := req.Identity.Get(ctx, &id); ds.HasError() {
		res.Diagnostics.Append(ds...)
		return
	}

	if id.ID.ValueString() != expectedID {
		// Unset ID
		res.Diagnostics.Append(res.Identity.Set(ctx, &FileResourceIdentity{ID: types.StringValue("")})...)
		return
	}
}

func (FileResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_base64": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sensitive_content": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_permission": schema.StringAttribute{
				Description: "Permissions to set for the output file",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0777"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4),
					stringvalidator.LengthAtLeast(3),
					&octalValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"directory_permission": schema.StringAttribute{
				Description: "Permissions to set for directories created",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0777"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4),
					stringvalidator.LengthAtLeast(3),
					&octalValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encrypted_regex": schema.StringAttribute{
				Description: "A regex pattern denoting the contents in the file to be encrypted",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"kms": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"arn": schema.StringAttribute{
						Description: "The ARN of the KMS key",
						Required:    true,
					},
					"profile": schema.StringAttribute{
						Description: "The AWS Profile to use when retrieving the key",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (FileResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

func (f FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var tfm FileResourceModel
	if ds := req.Config.Get(ctx, &tfm); ds.HasError() {
		resp.Diagnostics.Append(ds...)
		return
	}

	model := FileResourceAPIModel{
		Filename:            tfm.Filename.ValueString(),
		FilePermission:      tfm.FilePermission.ValueString(),
		DirectoryPermission: tfm.DirectoryPermission.ValueString(),
		SensitiveContent:    tfm.SensitiveContent.ValueString(),
		ContentBase64:       tfm.ContentBase64.ValueString(),
		Content:             tfm.Content.ValueString(),
		Source:              tfm.Source.ValueString(),
		EncryptedRegex:      tfm.EncryptedRegex.ValueString(),
	}

	if !tfm.Kms.IsNull() {
		var kmsConf KmsConf
		if ds := unmarshalKmsConf(ctx, tfm.Kms, &kmsConf); ds.HasError() {
			resp.Diagnostics.Append(ds...)
			return
		}
		model.Kms = &kmsConf
	}

	content, err := resourceLocalFileContent(model)
	if err != nil {
		resp.Diagnostics.AddError("base64 decode failure", err.Error())
		return
	}

	content, err = sopsEncrypt(ctx, model, content, &f.ProviderConfig)
	if err != nil {
		resp.Diagnostics.AddError("failed to encrypt", err.Error())
		return
	}

	//content = encrypt
	destinationDir := path.Dir(model.Filename)
	if _, err := os.Stat(destinationDir); err != nil {
		dirMode, _ := strconv.ParseInt(model.DirectoryPermission, 8, 64)
		if err := os.MkdirAll(destinationDir, os.FileMode(dirMode)); err != nil {
			resp.Diagnostics.AddError("failed to make directory for file", err.Error())
			return
		}
	}

	fileMode, _ := strconv.ParseInt(model.FilePermission, 8, 64)

	if err := os.WriteFile(model.Filename, content, os.FileMode(fileMode)); err != nil {
		resp.Diagnostics.AddError("failed to write file", err.Error())
		return
	}

	checksum := sha1.Sum(content)
	id := FileResourceIdentity{
		ID: types.StringValue(hex.EncodeToString(checksum[:])),
	}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &id)...)
}

func resourceLocalFileContent(f FileResourceAPIModel) ([]byte, error) {
	if f.SensitiveContent != "" {
		return []byte(f.SensitiveContent), nil
	}

	if f.ContentBase64 != "" {
		return base64.StdEncoding.DecodeString(f.ContentBase64)
	}

	if f.Source != "" {
		return os.ReadFile(f.Source)
	}

	return []byte(f.Content), nil
}

func sopsEncrypt(ctx context.Context, fr FileResourceAPIModel, content []byte, config *EncryptConfig) ([]byte, error) {
	inputStore := GetInputStore(fr.Filename)
	outputStore := GetOutputStore(fr.Filename)
	groups, err := KeyGroups(ctx, fr, "kms", config)
	if err != nil {
		return nil, err
	}

	encrypt, err := Encrypt(EncryptOpts{
		Cipher:         aes.NewCipher(),
		InputStore:     inputStore,
		OutputStore:    outputStore,
		InputPath:      fr.Filename,
		KeyServices:    LocalKeySvc(),
		EncryptedRegex: "kms",
		KeyGroups:      groups,
	}, content)

	if err != nil {
		return nil, err
	}

	return encrypt, nil
}

type octalValidator struct{}

func (v octalValidator) Description(_ context.Context) string {
	return "must be a file mode"
}

func (v octalValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v octalValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	fileMode, err := strconv.ParseInt(req.ConfigValue.ValueString(), 8, 64)
	if err != nil || fileMode > 0777 || fileMode < 0 {
		res.Diagnostics.AddAttributeError(req.Path, v.Description(ctx), "value must be a file mode")
	}
}
