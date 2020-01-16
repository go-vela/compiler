// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"context"
)

// key defines the key type for storing
// the compiler Engine in the context.
const key = "compiler"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the compiler Engine
// associated with this context.
func FromContext(c context.Context) Engine {
	// get compiler value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast compiler value to expected Engine type
	e, ok := v.(Engine)
	if !ok {
		return nil
	}

	return e
}

// ToContext adds the compiler Engine to this
// context if it supports the Setter interface.
func ToContext(c Setter, e Engine) {
	c.Set(key, e)
}
