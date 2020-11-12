package starlark

import (
	"reflect"
	"testing"

	"go.starlark.net/starlark"
)

func Test_toStarlark(t *testing.T) {
	dict := starlark.NewDict(16)
	err := dict.SetKey(starlark.String("foo"), starlark.String("bar"))
	if err != nil {
		t.Error(err)
	}
	a := make([]starlark.Value, 0)
	a = append(a, starlark.Value(starlark.String("foo")))
	a = append(a, starlark.Value(starlark.String("bar")))
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    starlark.Value
		wantErr bool
	}{
		{"string", args{value: "foo"}, starlark.String("foo"), false},
		{"byte array", args{value: []byte("array")}, starlark.String("array"), false},
		{"array", args{value: []string{"foo", "bar"}}, starlark.Tuple(a), false},
		{"bool", args{value: true}, starlark.Bool(true), false},
		{"float", args{value: 0.1}, starlark.Float(0.1), false},
		{"map", args{value: map[string]string{"foo": "bar"}}, dict, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toStarlark(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("toStarlark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toStarlark() got = %v, want %v", got, tt.want)
			}
		})
	}
}
