package symfony

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vairogs/cache-warmer/structs"
)

const (
	versionOption       = "--version"
	cacheWarmupArgument = "cache:warmup -q"
)

func CheckSymfonyConsole(config structs.Config) error {
	consoleFullPath := filepath.Join(config.SymfonyProjectDir, config.SymfonyConsolePath)
	if _, err := os.Stat(consoleFullPath); os.IsNotExist(err) {
		return fmt.Errorf("symfony console not found at %s", consoleFullPath)
	}
	return nil
}

// RunCommand executes a Symfony command with a given argument or option.
func RunCommand(config structs.Config, mainArgumentOrOption string) (string, error) {
	envOption := fmt.Sprintf("--env=%s", config.SymfonyEnv)
	consoleFullPath := filepath.Join(config.SymfonyProjectDir, config.SymfonyConsolePath)

	args := []string{mainArgumentOrOption, envOption}
	if !config.SymfonyDebug {
		args = append(args, "--no-debug")
	}

	out, err := exec.Command(consoleFullPath, args...).CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("symfony command failed: %s", exitErr.Error())
		}
		return "", fmt.Errorf("failed to execute Symfony command: %w", err)
	}

	return string(out), nil
}

// Version runs ./bin/console --version --env=dev
func Version(config structs.Config) (string, error) {
	return RunCommand(config, versionOption)
}

// CacheWarmup runs ./bin/console cache:warmup --env=dev
func CacheWarmup(config structs.Config) (string, error) {
	return RunCommand(config, cacheWarmupArgument)
}

func GetSymfonyProjectDir() (string, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return "", fmt.Errorf("no path provided")
	}

	path := args[0]

	if !filepath.IsAbs(path) {
		path = filepath.Join(execDir, path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	return path, nil
}
