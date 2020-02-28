// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

func TestNative_Compile_StagesPipeline(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	installEnv := environment(nil, nil, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})

	testEnv := environment(nil, nil, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	buildEnv := environment(nil, nil, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})

	dockerEnv := environment(nil, nil, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Template: false,
		},
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_init_init",
						Directory:   "/home//",
						Environment: environment(nil, nil, nil, nil),
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name: "clone",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_clone_clone",
						Directory:   "/home//",
						Environment: environment(nil, nil, nil, nil),
						Image:       "target/vela-git:v0.3.0",
						Name:        "clone",
						Number:      2,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "install",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_install_install",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: installEnv,
						Image:       "openjdk:latest",
						Name:        "install",
						Number:      3,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "test",
				Needs: []string{"install"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_test_test",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: testEnv,
						Image:       "openjdk:latest",
						Name:        "test",
						Number:      4,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "build",
				Needs: []string{"install"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_build_build",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: buildEnv,
						Image:       "openjdk:latest",
						Name:        "build",
						Number:      5,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "publish",
				Needs: []string{"build"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_publish_publish",
						Directory:   "/home//",
						Image:       "plugins/docker:18.09",
						Environment: dockerEnv,
						Name:        "publish",
						Number:      6,
						Pull:        true,
						Secrets: pipeline.StepSecretSlice{
							&pipeline.StepSecret{
								Source: "docker_username",
								Target: "registry_username",
							},
							&pipeline.StepSecret{
								Source: "docker_password",
								Target: "registry_password",
							},
						},
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/stages_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_StepsPipeline(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	installEnv := environment(nil, nil, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})

	testEnv := environment(nil, nil, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	buildEnv := environment(nil, nil, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})

	dockerEnv := environment(nil, nil, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_install",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/home//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Number:      3,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_test",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/home//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Number:      4,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_build",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/home//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: buildEnv,
				Image:       "openjdk:latest",
				Name:        "build",
				Number:      5,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_publish",
				Directory:   "/home//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "publish",
				Number:      6,
				Pull:        true,
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_StagesPipelineTemplate(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:name/contents/:path", func(c *gin.Context) {
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

	installEnv := environment(nil, nil, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})

	testEnv := environment(nil, nil, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	buildEnv := environment(nil, nil, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})

	dockerEnv := environment(nil, nil, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Template: false,
		},
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_init_init",
						Directory:   "/home//",
						Environment: environment(nil, nil, nil, nil),
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name: "clone",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_clone_clone",
						Directory:   "/home//",
						Environment: environment(nil, nil, nil, nil),
						Image:       "target/vela-git:v0.3.0",
						Name:        "clone",
						Number:      2,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "gradle",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_gradle_sample_install",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: installEnv,
						Image:       "openjdk:latest",
						Name:        "sample_install",
						Number:      3,
						Pull:        true,
					},
					&pipeline.Container{
						ID:          "__0_gradle_sample_test",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: testEnv,
						Image:       "openjdk:latest",
						Name:        "sample_test",
						Number:      4,
						Pull:        true,
					},
					&pipeline.Container{
						ID:          "__0_gradle_sample_build",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/home//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: buildEnv,
						Image:       "openjdk:latest",
						Name:        "sample_build",
						Number:      5,
						Pull:        true,
					},
				},
			},
			&pipeline.Stage{
				Name:  "publish",
				Needs: []string{"gradle"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_publish_publish",
						Directory:   "/home//",
						Image:       "plugins/docker:18.09",
						Environment: dockerEnv,
						Name:        "publish",
						Number:      6,
						Pull:        true,
						Secrets: pipeline.StepSecretSlice{
							&pipeline.StepSecret{
								Source: "docker_username",
								Target: "registry_username",
							},
							&pipeline.StepSecret{
								Source: "docker_password",
								Target: "registry_password",
							},
						},
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/stages_pipeline_template.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_StepsPipelineTemplate(t *testing.T) {
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

	installEnv := environment(nil, nil, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})

	testEnv := environment(nil, nil, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	buildEnv := environment(nil, nil, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})

	dockerEnv := environment(nil, nil, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_sample_install",
				Directory:   "/home//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "sample_install",
				Number:      3,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_sample_test",
				Directory:   "/home//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "sample_test",
				Number:      4,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_sample_build",
				Directory:   "/home//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: buildEnv,
				Image:       "openjdk:latest",
				Name:        "sample_build",
				Number:      5,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_docker",
				Directory:   "/home//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "docker",
				Number:      6,
				Pull:        true,
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_pipeline_template.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_InvalidType(t *testing.T) {
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

	gradleEnv := environment(nil, nil, nil, nil)
	gradleEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	gradleEnv["GRADLE_USER_HOME"] = ".gradle"

	dockerEnv := environment(nil, nil, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/home//",
				Environment: environment(nil, nil, nil, nil),
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        true,
			},
			&pipeline.Container{
				ID:          "step___0_docker",
				Directory:   "/home//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "docker",
				Number:      3,
				Pull:        true,
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/invalid_type.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_NoStepsorStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// run test
	yaml, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err == nil {
		t.Errorf("Compile should have returned err")
	}

	if got != nil {
		t.Errorf("Compile is %v, want %v", got, nil)
	}
}

func TestNative_Compile_StepsandStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_and_stages.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.Compile(yaml)
	if err == nil {
		t.Errorf("Compile should have returned err")
	}

	if got != nil {
		t.Errorf("Compile is %v, want %v", got, nil)
	}
}
