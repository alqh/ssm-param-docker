package compose

import (
	"os"
	"path/filepath"
)

// DefaultFileNames defines the Compose file names for auto-discovery (in order of preference)
var defaultFileNames = []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"}

// DefaultOverrideFileNames defines the Compose override file names for auto-discovery (in order of preference)
var defaultOverrideFileNames = []string{"compose.override.yml", "compose.override.yaml", "docker-compose.override.yml", "docker-compose.override.yaml"}

func FindComposeFiles(pwd string) []string {
	files := make([]string, 0, 2)

	for _, def := range defaultFileNames {
		f := filepath.Join(pwd, def)
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
			break
		}
	}

	for _, def := range defaultOverrideFileNames {
		f := filepath.Join(pwd, def)
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
			break
		}
	}

	return files
}
