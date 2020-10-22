package starlark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	types "github.com/go-vela/types/yaml"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"go.starlark.net/starlarkstruct"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, s *types.Step) (types.StepSlice, error) {
	// TODO: investigate way to not use tmp filesystem

	// setup filesystem
	file, err := ioutil.TempFile("/tmp", "sample")
	if err != nil {
		return nil, err
	}

	defer os.Remove(file.Name())

	_, err = file.Write([]byte(tmpl))
	if err != nil {
		return nil, err
	}

	config := new(types.Build)

	thread := &starlark.Thread{Name: "my thread"}
	globals, err := starlark.ExecFile(thread, file.Name(), nil, nil)
	if err != nil {
		return nil, err
	}

	mainVal, ok := globals["main"]
	if !ok {
		return nil, fmt.Errorf("no main function found")
	}
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, fmt.Errorf("main must be a function")
	}

	vars, err := userData(s.Template.Variables)
	if err != nil {
		return nil, err
	}

	args := starlark.Tuple([]starlark.Value{
		starlarkstruct.FromStringDict(
			starlark.String("context"), vars,
		),
	})

	// Call Starlark function from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			buf.WriteString("---\n")
			err = writeJSON(buf, item)
			if err != nil {
				return nil, err
			}
			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")
		err = writeJSON(buf, v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid return type (got a %s)", mainVal.Type())
	}

	fmt.Println(buf.String())

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buf.Bytes(), config)
	if err != nil {
		return types.StepSlice{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	fmt.Println(buf.String())

	return config.Steps, nil
}

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
		return fmt.Errorf("TypeError: value %s (type `%s') can't be converted to JSON.", v.String(), v.Type())
	}
	return nil
}

func userData(m map[string]interface{}) (starlark.StringDict, error) {
	dict := make(starlark.StringDict)

	for key, value := range m {
		val, err := toStarlark(value)
		if err != nil {
			return nil, err
		}
		dict[key] = val
	}

	return dict, nil
}

func toStarlark(vi interface{}) (starlark.Value, error) {
	if vi == nil {
		return starlark.None, nil
	}
	switch v := reflect.ValueOf(vi); v.Kind() {
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
		if b, ok := vi.([]byte); ok {
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

			d.SetKey(keyValue, kv)
		}
		return d, nil
	case reflect.Struct:
		ios, ok := vi.(intstr.IntOrString)
		if ok {
			switch ios.Type {
			case intstr.String:
				return starlark.String(ios.StrVal), nil
			case intstr.Int:
				return starlark.MakeInt(int(ios.IntVal)), nil
			}
		} else {
			data, err := json.Marshal(vi)
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
	return nil, fmt.Errorf("cannot convert %v to starlark", vi)
}
