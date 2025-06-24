package sops

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	wordwrap "github.com/mitchellh/go-wordwrap"

	mozillasops "github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/logging"

	//"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"

	//"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/version"
)

var log = logging.NewLogger("SOPS")

type EncryptOpts struct {
	Cipher            mozillasops.Cipher
	InputStore        mozillasops.Store
	OutputStore       mozillasops.Store
	InputPath         string
	KeyServices       []keyservice.KeyServiceClient
	UnencryptedSuffix string
	EncryptedSuffix   string
	UnencryptedRegex  string
	EncryptedRegex    string
	KeyGroups         []mozillasops.KeyGroup
	GroupThreshold    int
}

type fileAlreadyEncryptedError struct{}

func (err *fileAlreadyEncryptedError) Error() string {
	return "File already encrypted"
}

func (err *fileAlreadyEncryptedError) UserError() string {
	message := "The file you have provided contains a top-level entry called " +
		"'sops'. This is generally due to the file already being encrypted. " +
		"SOPS uses a top-level entry called 'sops' to store the metadata " +
		"required to decrypt the file. For this reason, SOPS can not " +
		"encrypt files that already contain such an entry.\n\n" +
		"If this is an unencrypted file, rename the 'sops' entry.\n\n" +
		"If this is an encrypted file and you want to edit it, use the " +
		"editor mode, for example: `sops my_file.yaml`"
	return wordwrap.WrapString(message, 75)
}

func ensureNoMetadata(branch mozillasops.TreeBranch) error {
	for _, b := range branch {
		if b.Key == "sops" {
			return &fileAlreadyEncryptedError{}
		}
	}
	return nil
}

func Encrypt(opts EncryptOpts, fileBytes []byte) (encryptedFile []byte, err error) {
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	if len(branches) == 0 {
		return nil, common.NewExitError(fmt.Sprintln("provided content was empty"), codes.CouldNotReadInputFile)
	}
	if err := ensureNoMetadata(branches[0]); err != nil {
		return nil, common.NewExitError(err, codes.FileAlreadyEncrypted)
	}
	path, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return nil, err
	}
	tree := mozillasops.Tree{
		Branches: branches,
		Metadata: mozillasops.Metadata{
			KeyGroups:         opts.KeyGroups,
			UnencryptedSuffix: opts.UnencryptedSuffix,
			EncryptedSuffix:   opts.EncryptedSuffix,
			UnencryptedRegex:  opts.UnencryptedRegex,
			EncryptedRegex:    opts.EncryptedRegex,
			Version:           version.Version,
			ShamirThreshold:   opts.GroupThreshold,
		},
		FilePath: path,
	}
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err = opts.OutputStore.EmitEncryptedFile(tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return
}

func LocalKeySvc() (svcs []keyservice.KeyServiceClient) {
	svcs = append(svcs, keyservice.NewLocalClient())
	return
}

func KeyGroups(ctx context.Context, fr fileResourceAPIModel, encType string, config *EncryptConfig) ([]mozillasops.KeyGroup, error) {
	var kmsKeys []keys.MasterKey
	switch encType {
	case "kms":
		if fr.Kms == nil || !fr.Kms.IsConfigured() {
			return nil, errors.New("kms is not configured")
		}

		for _, k := range kms.MasterKeysFromArnString(fr.Kms.ARN, nil, fr.Kms.Profile) {
			kmsKeys = append(kmsKeys, k)
		}
	}

	var group mozillasops.KeyGroup
	group = append(group, kmsKeys...)
	return []mozillasops.KeyGroup{group}, nil
}

type kmsConfigSchema struct {
	ARN     types.String `tfsdk:"arn"`
	Profile types.String `tfsdk:"profile"`
}

func unmarshalKmsConf(ctx context.Context, m types.Object, conf *KmsConf) diag.Diagnostics {
	var (
		ds       diag.Diagnostics
		tfSchema kmsConfigSchema
	)

	if m.IsNull() {
		// KMS is not configured
		return ds
	}

	if diags := m.As(ctx, &tfSchema, basetypes.ObjectAsOptions{}); diags.HasError() {
		ds.Append(diags...)
		return ds
	}

	// TODO: 'arn' only has to be specified on either the resource or the provider,
	// but this code does not distinguish between those two and runs on both - which means
	// that if the provider has no arn property but the resource does (or visa versa), this will fail
	if tfSchema.ARN.IsNull() {
		ds.AddAttributeError(tfpath.Root("arn"), "arn is not set", "arn is not set")
		return ds
	}

	conf.ARN = tfSchema.ARN.ValueString()
	// Profile is an optional value and is permitted to be an empty string if not specified.
	conf.Profile = tfSchema.Profile.ValueString()
	return ds
}
