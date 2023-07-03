package compose_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alqh/ssm-param-docker/compose"
)

func TestComposeSecrets_GetSecrets(t *testing.T) {
	t.Run("Extract secret names", func(t *testing.T) {
		secretConfigs, err := compose.ExtractSecrets([]string{"../test_data/secrets/docker-compose.yml"})
		require.NoError(t, err)

		names := make([]string, 0, len(secretConfigs))
		for _, sc := range secretConfigs {
			names = append(names, sc.Key)
		}
		sort.Strings(names)

		expected := []string{
			"domainname--servername--app_client_id",
			"domainname--servername--app_client_secret",
		}
		require.Equal(t, expected, names)
	})

	t.Run("Extract environment values", func(t *testing.T) {
		secretConfigs, err := compose.ExtractSecrets([]string{"../test_data/secrets/docker-compose.yml"})
		require.NoError(t, err)

		keyToEnv := make(map[string]string, 5)
		for _, s := range secretConfigs {
			keyToEnv[s.Key] = s.Environment
		}

		expected := map[string]string{
			"domainname--servername--app_client_id":     "MYAPP_CLIENT_ID",
			"domainname--servername--app_client_secret": "MYAPP_CLIENT_SECRET",
		}
		require.Equal(t, expected, keyToEnv)
	})
}
