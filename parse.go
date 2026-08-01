package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ComposeInfo holds the parsed metadata from a compose file.
type ComposeInfo struct {
	Name             string
	SourcePath       string
	ExternalNetworks []string
}

// composeFile is a minimal representation of a Docker Compose file.
type composeFile struct {
	Name     string                       `yaml:"name"`
	Networks map[string]composeNetworkDef `yaml:"networks"`
}

type composeNetworkDef struct {
	External bool `yaml:"external"`
}

// parseComposeFile reads a compose YAML and extracts name + external networks.
func parseComposeFile(path string) (*ComposeInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cf composeFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, err
	}

	name := cf.Name
	if name == "" {
		// Fallback: filename without extension.
		base := filepath.Base(path)
		name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	var extNets []string
	for netName, netDef := range cf.Networks {
		if netDef.External {
			extNets = append(extNets, netName)
		}
	}

	return &ComposeInfo{
		Name:             name,
		SourcePath:       path,
		ExternalNetworks: extNets,
	}, nil
}
