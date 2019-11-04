// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-vela/types/yaml"
)

// ScriptStages injects the script for each step in every stage in a yaml configuration.
func (c *Client) ScriptStages(s yaml.StageSlice) (yaml.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		// inject the scripts into the steps for the stage
		steps, err := c.ScriptSteps(stage.Steps)
		if err != nil {
			return nil, err
		}

		stage.Steps = steps
	}

	return s, nil
}

// ScriptSteps injects the script for each step in a yaml configuration.
func (c *Client) ScriptSteps(s yaml.StepSlice) (yaml.StepSlice, error) {
	// iterate through all steps
	for _, step := range s {
		// skip if no commands block for the step
		if len(step.Commands) == 0 {
			continue
		}

		// generate script from commands
		script := generateScriptPosix(step.Commands)

		// set the entrypoint for the step
		step.Entrypoint = []string{"/bin/sh", "-c"}

		// set the commands for the step
		step.Commands = []string{"echo $CARAVEL_BUILD_SCRIPT | base64 -d | /bin/sh -e"}

		// set the environment variables for the step
		step.Environment["CARAVEL_BUILD_SCRIPT"] = script
		step.Environment["HOME"] = "/root"
		step.Environment["SHELL"] = "/bin/sh"
	}

	return s, nil
}

// generateScriptPosix is a helper function that generates a build script
// for a linux container using the given commands
func generateScriptPosix(commands []string) string {
	var buf bytes.Buffer

	// iterate through each command provided
	for _, command := range commands {
		// safely escape entire command
		escaped := fmt.Sprintf("%q", command)

		// safely escape trace character
		escaped = strings.Replace(escaped, "$", `\$`, -1)

		// write escaped lines to buffer
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}

	// create build script with netrc and buffer information
	script := fmt.Sprintf(
		setupScript,
		buf.String(),
	)

	return base64.StdEncoding.EncodeToString([]byte(script))
}

// setupScript is a helper script this is added to the build to ensure
// a minimum set of environment variables are set correctly.
const setupScript = `
cat <<EOF > $HOME/.netrc
machine $CARAVEL_NETRC_MACHINE
login $CARAVEL_NETRC_USERNAME
password $CARAVEL_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
unset CARAVEL_NETRC_MACHINE
unset CARAVEL_NETRC_USERNAME
unset CARAVEL_NETRC_PASSWORD
unset CARAVEL_BUILD_SCRIPT
%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScript = `
echo $ %s
%s
`
