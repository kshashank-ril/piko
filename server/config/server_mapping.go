package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadServerMappings loads server mappings from a file or returns the in-memory mappings.
func (c *ServerMappingConfig) LoadServerMappings() (map[string]string, error) {
	if !c.Enabled {
		return nil, nil
	}

	// If we have in-memory mappings, use those
	if len(c.ServerMappings) > 0 {
		return c.ServerMappings, nil
	}

	// If we have a mapping file, load from it
	if c.MappingFile != "" {
		return c.loadFromFile(c.MappingFile)
	}

	return nil, nil
}

// loadFromFile loads server mappings from a YAML file.
func (c *ServerMappingConfig) loadFromFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open mapping file %s: %w", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read mapping file %s: %w", filename, err)
	}

	var mappings map[string]string
	if err := yaml.Unmarshal(data, &mappings); err != nil {
		return nil, fmt.Errorf("failed to parse mapping file %s: %w", filename, err)
	}

	return mappings, nil
}

// GetServerIP returns the server IP for a given endpoint ID.
func (c *ServerMappingConfig) GetServerIP(endpointID string) (string, bool) {
	if !c.Enabled {
		return "", false
	}

	mappings, err := c.LoadServerMappings()
	if err != nil {
		return "", false
	}

	// Split endpoint ID by underscore and take the first part
	parts := strings.Split(endpointID, "_")
	if len(parts) == 0 {
		return "", false
	}

	prefix := parts[0]
	serverIP, exists := mappings[prefix]
	return serverIP, exists
}
