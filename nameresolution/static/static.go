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
	"strconv"

	"github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/dapr/pkg/logger"
)

type resolver struct {
	logger logger.Logger
}

type StaticEntry struct {
	Address string `json:"host"`
	Port int       `json:"port"`
}

// NewResolver creates static name resolver.
func NewResolver(logger logger.Logger) nameresolution.Resolver {
	return &resolver{logger: logger}
}

// Init initializes static name resolver.
func (k *resolver) Init(metadata nameresolution.Metadata) error {
	var port, err = strconv.ParseInt(metadata.Properties[nameresolution.MDNSInstancePort], 10, 32)
	if (err != nil) {
		return err
	}

	var entry = StaticEntry{
		Address: metadata.Properties[nameresolution.MDNSInstanceAddress],
		Port: int(port),
	}

	return k.writeConfigFileForApp(metadata.Properties[nameresolution.MDNSInstanceName], entry)
}

// ResolveID resolves name to static address.
func (k *resolver) ResolveID(req nameresolution.ResolveRequest) (string, error) {
	var config, err = k.readConfigFileForApp(req.ID)
	if (err != nil) {
		return "", err
	}

	var address = fmt.Sprintf("%s:%d", config.Address, config.Port)
	k.logger.Infof("Resolved address for ID %s to %s", req.ID, address)
	return address, nil
}

func (k *resolver) getFileNameForApp(id string) string {
	return fmt.Sprintf("%s.json", id)
}

func (k *resolver) readConfigFileForApp(id string) (*StaticEntry, error) {
	var fileName = k.getFileNameForApp(id)

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

func (k *resolver) writeConfigFileForApp(id string, entry StaticEntry) error {
	var byteValue, err = json.Marshal(entry)
	if (err != nil) {
		return err
	}

	var fileName = k.getFileNameForApp(id)

	err = ioutil.WriteFile(fileName, byteValue, 0644 /* TODO: Is this permission reasonable? */)
	if err != nil {
		return err
	}

	return nil
}