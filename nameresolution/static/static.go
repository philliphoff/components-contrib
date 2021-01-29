// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package static

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/dapr/pkg/logger"
)

type resolver struct {
	logger logger.Logger
}

type StaticEntry struct {
	Host string `json:"host"`
}

// NewResolver creates static name resolver.
func NewResolver(logger logger.Logger) nameresolution.Resolver {
	return &resolver{logger: logger}
}

// Init initializes static name resolver.
func (k *resolver) Init(metadata nameresolution.Metadata) error {
	return nil
}

// ResolveID resolves name to static address.
func (k *resolver) ResolveID(req nameresolution.ResolveRequest) (string, error) {
	var config, err = k.readConfigFileForApp(req.ID)
	if (err != nil) {
		return "", err
	}

	var address = fmt.Sprintf("%s:%d", config.Host, req.Port)
	k.logger.Infof("Resolved address for ID %s to %s", req.ID, address)
	return address, nil
}

func (k *resolver) readConfigFileForApp(id string) (*StaticEntry, error) {
	var fileName = fmt.Sprintf("%s.json", id)

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var jsonConfig StaticEntry
	err = json.Unmarshal(byteValue, &jsonConfig)
	if err != nil {
		return nil, err
	}

	return &jsonConfig, nil
}