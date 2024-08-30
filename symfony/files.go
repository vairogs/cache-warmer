package symfony

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vairogs/cache-warmer/structs"
)

// FindFiles searches for files in the specified root directory and its subdirectories.
// It excludes the specified directories and handles vendor directories based on the vendorWatch flag.
// It returns a list of file paths and an error if any occurred.
func FindFiles(root string, excludedDirs []string, vendorWatch bool, vendorList []string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip excluded directories
			for _, excludedDir := range excludedDirs {
				if strings.Contains(path, excludedDir) {
					return filepath.SkipDir
				}
			}

			// Handle vendor directory based on the vendorWatch flag
			if strings.Contains(path, "vendor") && !vendorWatch {
				return filepath.SkipDir
			}

			if vendorWatch {
				insideVendor := false
				for _, vendor := range vendorList {
					if strings.HasPrefix(path, filepath.Join(root, vendor)) {
						insideVendor = true
						break
					}
				}
				if !insideVendor {
					return filepath.SkipDir
				}
			}
		} else {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// GetWatchMap returns a map containing the files to watch and their corresponding last modified timestamps.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values for the application.
// It calls the `GetFilesToWatch` function to retrieve the files to watch.
// For each file, it retrieves the file's stats using `os.Stat`, and adds the file path and last modified timestamp
// to the `watchMap`. If any error occurs during the process, it returns the error.
// It returns the `watchMap` containing the files to watch and their corresponding timestamps and nil error on success.
// Otherwise, it returns nil map and the encountered error.
// The `watchMap` can be used to compare with the existing files being watched to detect any changes.
//
// Note that `GetWatchMap` does not handle removing files from the `watchMap` when they are no longer being watched.
// This responsibility falls on the caller of this function.
func GetWatchMap(config structs.Config) (map[string]string, error) {
	watchMap := make(map[string]string)
	filesToWatch, err := GetFilesToWatch(config)
	if err != nil {
		return nil, err
	}

	for _, file := range filesToWatch {
		stats, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("can't get stats for the \"%s\" file, check the project permissions or if a new file was created: %v", file, err)
		}
		watchMap[file] = stats.ModTime().String()
	}

	return watchMap, nil
}

func GetFilesToWatch(config structs.Config) ([]string, error) {
	var filesToWatch []string

	// Set up excluded directories
	var excludedDirs = config.ExcludeDirs
	if !config.VendorWatch {
		excludedDirs = append(excludedDirs, "vendor")
	}

	// Include general files like .env*
	generalFiles, err := GetFilesFromPath(config, ".env*")
	if err != nil {
		return nil, err
	}
	filesToWatch = append(filesToWatch, generalFiles...)

	// Directories to watch
	symfonyDirs := map[string]string{
		config.SymfonyConfigDir:       "config",
		config.SymfonySrcDir:          "src",
		config.SymfonyTemplatesDir:    "templates",
		config.SymfonyTranslationsDir: "translations",
		config.MigrationsDir:          "migrations",
	}

	// Watch all files in the specified directories, regardless of their extensions
	for dir := range symfonyDirs {
		files, err := FindFiles(dir, excludedDirs, config.VendorWatch, config.VendorList)
		if err != nil {
			return nil, err
		}
		filesToWatch = append(filesToWatch, files...)
	}

	// If VendorWatch is enabled, watch specific vendor directories
	if config.VendorWatch {
		for _, vendor := range config.VendorList {
			vendorPath := filepath.Join(config.SymfonyVendorDir, vendor)
			vendorFiles, err := FindFiles(vendorPath, excludedDirs, true, config.VendorList)
			if err != nil {
				return nil, err
			}
			filesToWatch = append(filesToWatch, vendorFiles...)
		}
	}

	return filesToWatch, nil
}

// GetFilesFromPath retrieves a list of files from the specified path based on the provided configuration.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values.
// The function uses the filepath package to glob files based on the path and the provided glob pattern.
// If an error occurs during the globbing process, the function returns an error.
// Otherwise, it returns the list of files matching the pattern and nil error.
func GetFilesFromPath(config structs.Config, glob string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(config.SymfonyProjectDir, glob))
	if err != nil {
		return nil, fmt.Errorf("error while globbing files: %v", err)
	}

	return files, nil
}
