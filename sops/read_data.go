package sops

import (
	"encoding/json"
	"fmt"

	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v2"

	"github.com/lokkersp/terraform-provider-sops/sops/internal/dotenv"
	"github.com/lokkersp/terraform-provider-sops/sops/internal/ini"
)

// readData consolidates the logic of extracting the from the various input methods and setting it on the ResourceData
func readData(content []byte, format string, model *FileDataSourceModel) error {
	cleartext, err := decrypt.Data(content, format)
	if err != nil {
		return err
	}

	// Set output attribute for raw content
	model.Raw = string(cleartext)

	// Set output attribute for content as a map (only for json and yaml)
	var data map[string]any
	switch format {
	case "json":
		err = json.Unmarshal(cleartext, &data)
	case "yaml":
		err = yaml.Unmarshal(cleartext, &data)
	case "dotenv":
		err = dotenv.Unmarshal(cleartext, &data)
	case "ini":
		err = ini.Unmarshal(cleartext, &data)
	}

	if err != nil {
		return err
	}

	model.ID = "-"
	model.Data = flatten(data)
	return nil
}

type readDataKeyModel struct {
	ID   string
	Data string
	Raw  string
	Map  map[string]string
	Yaml string
}

// readData consolidates the logic of extracting the from the various input methods and setting it on the ResourceData
func readDataKey(content []byte, format string, key string, d *readDataKeyModel) error {
	cleartext, err := decrypt.Data(content, format)
	if err != nil {
		return fmt.Errorf("fail to decrypt,format is %s:%s", format, err)
	}

	d.Raw = string(cleartext)
	// Set output attribute for content as a map (only for json and yaml)
	var data map[string]interface{}
	switch format {
	case "json":
		err = json.Unmarshal(cleartext, &data)
	case "yaml":
		err = yaml.Unmarshal(cleartext, &data)
	case "dotenv":
		err = dotenv.Unmarshal(cleartext, &data)
	case "ini":
		err = ini.Unmarshal(cleartext, &data)
	}
	if err != nil {
		return fmt.Errorf("evaluated format is %s:%s", err, format)
	}

	d.Data = flatten(data)[key]

	err, value := flattenFromKey(data, key)
	out, err := yaml.Marshal(map[string]interface{}{key: data[key]})
	if err != nil {
		return err
	}

	d.Map = value
	d.Yaml = string(out)
	d.ID = "-"
	return nil
}
