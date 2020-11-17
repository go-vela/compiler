// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func TestNative_ExpandStages(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	tmpls := map[string]*yaml.Template{
		"gradle": {
			Name:   "gradle",
			Source: "github.example.com/foo/bar/template.yml",
			Type:   "github",
		},
	}

	stages := yaml.StageSlice{
		&yaml.Stage{
			Name: "foo",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Name: "sample",
					Template: yaml.StepTemplate{
						Name: "gradle",
						Variables: map[string]interface{}{
							"image":       "openjdk:latest",
							"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
							"pull_policy": "pull: true",
						},
					},
				},
			},
		},
	}

	want := yaml.StageSlice{
		&yaml.Stage{
			Name: "foo",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Commands: []string{"./gradlew downloadDependencies"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_install",
					Pull:  "always",
				},
				&yaml.Step{
					Commands: []string{"./gradlew check"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_test",
					Pull:  "always",
				},
				&yaml.Step{
					Commands: []string{"./gradlew build"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_build",
					Pull:  "always",
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	got, err := compiler.ExpandStages(stages, tmpls)
	if err != nil {
		t.Errorf("ExpandStages returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExpandStages is %v, want %v", got, want)
	}
}

func TestNative_ExpandSteps(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	tmpls := map[string]*yaml.Template{
		"gradle": {
			Name:   "gradle",
			Source: "github.example.com/foo/bar/template.yml",
			Type:   "github",
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "gradle",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_build",
			Pull:  "always",
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	got, err := compiler.ExpandSteps(steps, tmpls)
	if err != nil {
		t.Errorf("ExpandSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExpandSteps is %v, want %v", got, want)
	}
}

func TestNative_ExpandStepsStarlark(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template-starlark.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	tmpls := map[string]*yaml.Template{
		"go": {
			Name:   "go",
			Source: "github.example.com/foo/bar/template.star",
			Format: "starlark",
			Type:   "github",
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name:      "go",
				Variables: map[string]interface{}{},
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"go build", "go test"},
			Image:    "golang:latest",
			Name:     "sample_build",
			Pull:     "not_present",
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	got, err := compiler.ExpandSteps(steps, tmpls)
	if err != nil {
		t.Errorf("ExpandSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
		}
		t.Errorf("ExpandSteps is %v, want %v", got, want)
	}
}

func TestNative_mapFromTemplates(t *testing.T) {
	// setup types
	str := "foo"

	tmpl := []*yaml.Template{
		{
			Name:   str,
			Source: str,
			Type:   str,
		},
	}

	want := map[string]*yaml.Template{
		str: {
			Name:   str,
			Source: str,
			Type:   str,
		},
	}

	// run test
	got := mapFromTemplates(tmpl)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromTemplates is %v, want %v", got, want)
	}
}
