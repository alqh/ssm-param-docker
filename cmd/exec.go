package cmd

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/alqh/ssm-param-docker/aws"
	"github.com/alqh/ssm-param-docker/compose"
)

func ExecCLI() *cli.Command {
	return &cli.Command{
		Name:   "exec",
		Usage:  "Pull secrets defined in docker compose definition and exec the command",
		Action: ExecCmd,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "env",
				Usage: "ENV_KEY=ENV_VALUE to override the environment variables",
			},
			&cli.StringFlag{
				Name:        "param-path",
				Usage:       "The hierarchy path of where the parameters are stored. The compose file property key can also contain the hierarchy via -- separator",
				DefaultText: "/",
			},
			&cli.StringSliceFlag{
				Name:  "file",
				Usage: "Overriding the docker compose files to use",
			},
			&cli.StringFlag{
				Name:  "aws-region",
				Usage: "AWS region to use to retrieve the ssm parameter",
			},
			&cli.StringFlag{
				Name:  "aws-profile",
				Usage: "AWS profile to use to retrieve the ssm parameter",
			},
		},
	}
}

func ExecCmd(cCtx *cli.Context) error {
	// Figure out the compose file paths.
	filesArg := cCtx.StringSlice("file")
	files := filesArg
	if len(filesArg) == 0 {
		// Scan for compose files.
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "unable to get current working directory")
		}
		files = compose.FindComposeFiles(cwd)
	} else {
		// Validate files passed in args.
		for _, f := range files {
			if _, err := os.Stat(f); err != nil {
				return errors.Wrap(err, "file not valid")
			}
		}
	}

	// Extract secrets config in the compose files.
	secrets, err := compose.ExtractSecrets(files)
	if err != nil {
		return errors.Wrap(err, "unable to extract secrets config from compose file")
	}

	// Pull down parameter store.
	cfgOpts := make([]func(options *config.LoadOptions) error, 0, 2)
	if cCtx.String("aws-region") != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(cCtx.String("aws-region")))
	}
	if cCtx.String("aws-profile") != "" {
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(cCtx.String("aws-profile")))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
	if err != nil {
		return errors.Wrap(err, "error creating aws ssm client")
	}
	client := ssm.NewFromConfig(cfg)
	paramStore := aws.NewAWSParamStore(client)

	paramPath := cCtx.String("param-path")
	paramKeyToSecret := make(map[string]compose.SecretConfig, len(secrets))
	paramKeys := make([]string, 0, len(secrets))
	for _, s := range secrets {
		// Convert delimiter to a path.
		key := strings.ReplaceAll(s.Key, "--", "/")
		paramKeys = append(paramKeys, key)
		paramKeyToSecret[key] = s
	}

	paramVals, err := paramStore.Parameters(cCtx.Context, paramPath, paramKeys)
	if err != nil {
		return errors.Wrapf(err, "failed to extract parameters %v within hierarchy %s", paramKeys, paramPath)
	}

	envParamToValue := make(map[string]string, len(secrets))
	for paramKey, paramVal := range paramVals {
		secret := paramKeyToSecret[paramKey]
		envParamToValue[secret.Environment] = paramVal
	}

	// Exec command with the values.
	cmd := cCtx.Args().Get(0)
	if cmd == "" {
		log.Printf("No command to run")
		return nil
	}

	if err := RunCmd(cmd, cCtx.Args().Tail(), cCtx.StringSlice("env"), envParamToValue); err != nil {
		return errors.Wrap(err, "error executing command")
	}
	return nil
}
