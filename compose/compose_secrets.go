package compose

import (
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

type ExternalSecret struct {
	External bool
	Name     string
}

type SecretConfig struct {
	Key         string
	Environment string
	External    ExternalSecret
}

func ExtractSecrets(composeFiles []string) ([]SecretConfig, error) {
	composeCfg := types.ConfigDetails{
		ConfigFiles: make([]types.ConfigFile, 0, len(composeFiles)),
	}

	for _, f := range composeFiles {
		composeCfg.ConfigFiles = append(composeCfg.ConfigFiles, types.ConfigFile{Filename: f})
	}

	cfg, err := loader.Load(composeCfg)
	if err != nil {
		return nil, err
	}

	secrets := make([]SecretConfig, 0, len(cfg.Secrets))
	for k, s := range cfg.Secrets {
		sc := SecretConfig{
			Key:         k,
			Environment: s.Environment,
			External: ExternalSecret{
				External: s.External.External,
				Name:     s.External.Name,
			},
		}
		secrets = append(secrets, sc)
	}

	return secrets, nil
}
