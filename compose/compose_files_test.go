package compose_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alqh/ssm-param-docker/compose"
)

func TestComposeFiles_FindFiles(t *testing.T) {
	t.Run("Prefer compose.yml over docker_compose.yml", func(t *testing.T) {
		f := compose.FindComposeFiles("../test_data/compose_files/compose")

		require.Equal(t, "../test_data/compose_files/compose/compose.yml", f[0])
		require.NotContains(t, f, "../test_data/compose_files/compose/docker-compose.yml")
	})

	t.Run("Finds docker_compose.yml", func(t *testing.T) {
		f := compose.FindComposeFiles("../test_data/compose_files/docker_compose")

		require.Equal(t, "../test_data/compose_files/docker_compose/docker-compose.yml", f[0])
		require.NotContains(t, f, "../test_data/compose_files/docker_compose/compose.yml")
	})

	t.Run("Allow compose override", func(t *testing.T) {
		f := compose.FindComposeFiles("../test_data/compose_files/compose")

		require.Equal(t, "../test_data/compose_files/compose/compose.override.yml", f[1])
	})

	t.Run("Allow docker_compose override", func(t *testing.T) {
		f := compose.FindComposeFiles("../test_data/compose_files/docker_compose")

		require.Equal(t, "../test_data/compose_files/docker_compose/docker-compose.yml", f[0])
		require.NotContains(t, f, "../test_data/compose_files/docker_compose/compose.yml")
	})
}
