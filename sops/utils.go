package sops

import (
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/config"
)

var defaultStoreConfig = config.NewStoresConfig()

func GetInputStore(filename string) common.Store {
	return common.DefaultStoreForPathOrFormat(defaultStoreConfig, filename, "file")
}
func GetOutputStore(filename string) common.Store {
	return common.DefaultStoreForPathOrFormat(defaultStoreConfig, filename, "file")
}
