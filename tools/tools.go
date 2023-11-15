// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	// copywrite header generation
	_ "github.com/hashicorp/copywrite"
)
