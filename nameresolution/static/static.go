// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package static

import (
	"fmt"

	"github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/dapr/pkg/logger"
)

type resolver struct {
	logger logger.Logger
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
	var address = fmt.Sprintf("localhost:%d", req.Port)
	k.logger.Infof("Resolved address for ID %s to %s", req.ID, address)
	return address, nil
}
