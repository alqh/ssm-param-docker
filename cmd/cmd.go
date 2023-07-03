package cmd

import (
	"log"
	"os"
	"os/exec"
)

type environ []string

func RunCmd(cmd string, args []string, overrideEnvs []string, envParamToValue map[string]string) error {

	existingEnv := os.Environ()
	envs := environ(existingEnv)

	for envK, envV := range envParamToValue {
		envs = append(envs, envK+"="+envV)
	}

	if len(overrideEnvs) > 0 {
		for _, v := range overrideEnvs {
			envs = append(envs, v)
		}
	}

	c := exec.Command(cmd, args...)
	c.Env = envs

	out, err := c.Output()
	if err != nil {
		return err
	}
	log.Println(string(out))
	return nil
}
