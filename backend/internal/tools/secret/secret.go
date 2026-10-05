package secret

import (
	"fmt"
	"os"
	"strings"
)

// GetSecret lit la variable d'environnement envVarName pour obtenir le chemin du fichier,
// puis retourne le contenu nettoyé (sans espaces/sauts de ligne superflus).
func GetSecret(envVarName string) (string, error) {
	filePath := os.Getenv(envVarName)
	if filePath == "" {
		return "", fmt.Errorf("Env variable %s is not defined", envVarName)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Cannot read the secret path for env variable %s : %w", envVarName, err)
	}

	// deletes any whitespace
	return strings.TrimSpace(string(content)), nil
}