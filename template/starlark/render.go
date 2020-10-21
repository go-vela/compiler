package starlark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	types "github.com/go-vela/types/yaml"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"

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

	args := starlark.Tuple([]starlark.Value{
		starlarkstruct.FromStringDict(
			starlark.String("context"), starlark.StringDict{},
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
			data, _ := json.Marshal(s)
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
