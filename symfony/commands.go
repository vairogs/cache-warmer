package symfony

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vairogs/cache-warmer/structs"
)

const (
	versionOption       = "--version"
	cacheWarmupArgument = "cache:warmup"
	cacheClearArgument  = "cache:clear"
	cachePoolArgument   = "cache:pool:clear"
)

// CheckSymfonyConsole checks if the Symfony console exists at the specified path in the given configuration.
// If the console does not exist, an error is returned.
// The function takes a configuration object as a parameter and uses the Symfony project directory and the relative
// path to the Symfony console from the configuration to generate the full path to the console.
// It then checks if the console file exists at that path and returns an error if it does not.
// Otherwise, it returns nil.
func CheckSymfonyConsole(config structs.Config) error {
	consoleFullPath := filepath.Join(config.SymfonyProjectDir, config.SymfonyConsolePath)
	if _, err := os.Stat(consoleFullPath); os.IsNotExist(err) {
		return fmt.Errorf("symfony console not found at %s", consoleFullPath)
	}

	return nil
}

// RunCommand executes a Symfony console command with the provided configuration and main argument or option.
// It returns the combined output of the command and an error, if any.
// The config parameter is an instance of the Config struct, which holds the necessary parameters for the application.
// The mainArgumentOrOption parameter is the main argument or option to be passed to the Symfony console command.
// The function constructs the command with the appropriate arguments based on the config and executes it using the exec package.
// If the command fails, it returns an error with a relevant error message.
// The return value is the output of the command as a string and any error encountered during execution.
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

// Version executes a Symfony console command with the provided configuration and the "--version" option.
// It returns the combined output of the command as a string and any error encountered during execution.
// The config parameter is an instance of the Config struct, which holds the necessary parameters for the application.
// The function calls the RunCommand function passing the config and versionOption as arguments.
// The return value is the output of the command as a string and any error encountered.
func Version(config structs.Config) (string, error) {
	return RunCommand(config, versionOption)
}

// CacheWarmup warms up the cache based on the provided configuration.
// If the config.ClearCache flag is set to true, it clears the cache using the cache:clear command.
// If the config.ForceClearCache flag is set to true, it removes the cache directory using the rm -rf command.
// If the config.PoolsProvided flag is set to true, it clears the specified cache pools using the cache:pool:clear command.
// Finally, it warms up the cache using the cache:warmup command.
// The function returns the output of the cache:warmup command as a string and any error encountered during execution.
func CacheWarmup(config structs.Config) (string, error) {
	if config.ClearCache {
		_, err := RunCommand(config, cacheClearArgument)
		if err != nil {
			return "", fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	if config.ForceClearCache {
		err := RemoveCache(config)
		if err != nil {
			return "", fmt.Errorf("failed to remove cache: %w", err)
		}
	}

	if config.PoolsProvided {
		for _, pool := range config.Pools {
			_, err := RunCommand(config, cachePoolArgument+" "+pool)
			if err != nil {
				return "", fmt.Errorf("failed to clear pool: %w", err)
			}
		}
	}

	return RunCommand(config, cacheWarmupArgument)
}

// RemoveCache removes the cache directory based on the provided configuration.
// It takes a Config object as its parameter, which holds the necessary parameters
// for the application. The cache directory is removed using the "rm -rf" command.
// If the cache directory is not within the project directory, an error is returned.
// If the removal of the cache directory fails, an error is returned.
// The function returns an error if any error occurs, otherwise it returns nil.
func RemoveCache(config structs.Config) error {
	projectDir := config.SymfonyProjectDir
	cacheDir := filepath.Join(projectDir, "var", "cache")

	if !strings.HasPrefix(cacheDir, projectDir) {
		return fmt.Errorf("invalid projectDir: %s is not within the root directory", cacheDir)
	}

	cmd := exec.Command("rm", "-rf", cacheDir)
	cmd.Dir = projectDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove cache directory: %v", err)
	}

	return nil
}

// GetSymfonyProjectDir returns the path to the Symfony project directory based on the command line arguments.
// It retrieves the current working directory, parses the command line arguments, and checks for the existence of a provided path.
// If the provided path is relative, it joins it with the current working directory.
// If the provided path does not exist, it returns an error.
// Otherwise, it returns the path to the Symfony project directory and nil error.
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
