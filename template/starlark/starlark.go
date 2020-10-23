// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// toStarlark takes an value as an interface an
// will return the comparable Starlark type.
//
// This code is under copyright (full attribution in NOTICE) and is from:
// https://github.com/wonderix/shalm/blob/899b8f7787883d40619eefcc39bd12f42a09b5e7/pkg/shalm/convert.go#L14-L85
// nolint /// @lll ignore line length of the link
func toStarlark(value interface{}) (starlark.Value, error) {
	logrus.Tracef("converting %v to starklark type", value)

	if value == nil {
		return starlark.None, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		return starlark.String(v.String()), nil
	case reflect.Bool:
		return starlark.Bool(v.Bool()), nil
	case reflect.Int:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int16:
		return starlark.MakeInt64(v.Int()), nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint16:
		return starlark.MakeUint64(v.Uint()), nil
	case reflect.Float32:
		return starlark.Float(v.Float()), nil
	case reflect.Float64:
		return starlark.Float(v.Float()), nil
	case reflect.Slice:
		if b, ok := value.([]byte); ok {
			return starlark.String(string(b)), nil
		}

		a := make([]starlark.Value, 0)

		for i := 0; i < v.Len(); i++ {
			val, err := toStarlark(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			a = append(a, val)
		}

		return starlark.Tuple(a), nil
	case reflect.Ptr:
		val, err := toStarlark(v.Elem().Interface())
		if err != nil {
			return nil, err
		}

		return val, nil
	case reflect.Map:
		d := starlark.NewDict(16)

		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			keyValue, err := toStarlark(key.Interface())
			if err != nil {
				return nil, err
			}

			kv, err := toStarlark(strct.Interface())
			if err != nil {
				return nil, err
			}

			err = d.SetKey(keyValue, kv)
			if err != nil {
				return nil, err
			}
		}

		return d, nil
	case reflect.Struct:
		ios, ok := value.(intstr.IntOrString)
		if ok {
			switch ios.Type {
			case intstr.String:
				return starlark.String(ios.StrVal), nil
			case intstr.Int:
				return starlark.MakeInt(int(ios.IntVal)), nil
			}
		} else {
			data, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			var m map[string]interface{}
			err = json.Unmarshal(data, &m)
			if err != nil {
				return nil, err
			}

			return toStarlark(m)
		}
	}

	return nil, fmt.Errorf("unable to convert %v to starlark", value)
}

// writeJSON takes an starlark input and return the valid JSON
// for the specific type.
//
// This code is under copyright (full attribution in NOTICE) and is from:
// https://github.com/drone/drone-cli/blob/master/drone/starlark/starlark.go#L214-L274
func writeJSON(out *bytes.Buffer, v starlark.Value) error {
	if marshaler, ok := v.(json.Marshaler); ok {
		jsonData, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}
		out.Write(jsonData)
		return nil
	}

	switch v := v.(type) {
	case starlark.NoneType:
		out.WriteString("null")
	case starlark.Bool:
		fmt.Fprintf(out, "%t", v)
	case starlark.Int:
		out.WriteString(v.String())
	case starlark.Float:
		fmt.Fprintf(out, "%g", v)
	case starlark.String:
		s := string(v)
		if goQuoteIsSafe(s) {
			fmt.Fprintf(out, "%q", s)
		} else {
			// vanishingly rare for text strings
			data, err := json.Marshal(s)
			if err != nil {
				logrus.Error(err)
			}
			out.Write(data)
		}
	case starlark.Indexable: // Tuple, List
		out.WriteByte('[')
		for i, n := 0, starlark.Len(v); i < n; i++ {
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, v.Index(i)); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case *starlark.Dict:
		out.WriteByte('{')
		for i, itemPair := range v.Items() {
			key := itemPair[0]
			value := itemPair[1]
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, key); err != nil {
				return err
			}
			out.WriteString(": ")
			if err := writeJSON(out, value); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("unable to convert value %s (type `%s') can't be converted to json.", v.String(), v.Type())
	}
	return nil
}

// goQuoteIsSafe takes a string and checks if is safely quoted
//
// This code is under copyright (full attribution in NOTICE) and is from:
// https://github.com/drone/drone-cli/blob/master/drone/starlark/starlark.go#L276-L285
func goQuoteIsSafe(s string) bool {
	for _, r := range s {
		// JSON doesn't like Go's \xHH escapes for ASCII control codes,
		// nor its \UHHHHHHHH escapes for runes >16 bits.
		if r < 0x20 || r >= 0x10000 {
			return false
		}
	}
	return true
}
