// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	types "github.com/go-vela/types/yaml"

	"github.com/buildkite/yaml"
)

// Parse converts an object to a yaml configuration.
func (c *client) Parse(v interface{}) (*types.Build, error) {
	switch v := v.(type) {
	case []byte:
		return ParseBytes(v)
	case *os.File:
		return ParseFile(v)
	case io.Reader:
		return ParseReader(v)
	case string:
		// check if string is path to file
		_, err := os.Stat(v)
		if err == nil {
			// parse string as path to yaml configuration
			return ParsePath(v)
		}

		// parse string as yaml configuration
		return ParseString(v)
	default:
		return nil, fmt.Errorf("unable to parse yaml: unrecognized type %T", v)
	}
}

// ParseBytes converts a byte slice to a yaml configuration.
func ParseBytes(b []byte) (*types.Build, error) {
	config := new(types.Build)

	// unmarshal the bytes into the yaml configuration
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	fmt.Println("METADATA 0: ", config.Metadata.Clone)

	return config, nil
}

// ParseFile converts an os.File into a yaml configuration.
func ParseFile(f *os.File) (*types.Build, error) {
	return ParseReader(f)
}

// ParsePath converts a file path into a yaml configuration.
func ParsePath(p string) (*types.Build, error) {
	// open the file for reading
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("unable to open yaml file %s: %v", p, err)
	}

	defer f.Close()

	return ParseReader(f)
}

// ParseReader converts an io.Reader into a yaml configuration.
func ParseReader(r io.Reader) (*types.Build, error) {
	// read all the bytes from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read bytes for yaml: %v", err)
	}

	return ParseBytes(b)
}

// ParseString converts a string into a yaml configuration.
func ParseString(s string) (*types.Build, error) {
	return ParseBytes([]byte(s))
}
