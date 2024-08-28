package symfony

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vairogs/cache-warmer/structs"
)

func findFiles(root string, extensions []string, excludedDirs []string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			for _, excludedDir := range excludedDirs {
				if strings.Contains(path, excludedDir) {
					return nil
				}
			}
			for _, ext := range extensions {
				if filepath.Ext(d.Name()) == ext {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})
	return files, err
}

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

	fileExtensions := []string{
		".yaml", ".yml", ".php", ".twig", ".xlf", ".csv", ".json", ".dat", ".res", ".mo", ".po", ".qt",
	}

	generalFiles, err := getFilesFromPath(config, ".env*")
	if err != nil {
		return nil, err
	}
	filesToWatch = append(filesToWatch, generalFiles...)

	symfonyDirs := map[string]string{
		config.SymfonyConfigDir:       "config",
		config.SymfonySrcDir:          "src",
		config.SymfonyTemplatesDir:    "templates",
		config.SymfonyTranslationsDir: "translations",
	}

	var excludeDirs = config.ExcludeDirs
	if !config.VendorWatch {
		excludeDirs = append(excludeDirs, "vendor")
	}

	for dir := range symfonyDirs {
		files, err := findFiles(dir, fileExtensions, excludeDirs)
		if err != nil {
			return nil, err
		}
		filesToWatch = append(filesToWatch, files...)
	}

	if config.VendorWatch {
		for _, vendor := range config.VendorList {
			vendorPath := filepath.Join(config.SymfonyVendorDir, vendor)
			vendorFiles, err := findFiles(vendorPath, fileExtensions, excludeDirs)
			if err != nil {
				return nil, err
			}
			filesToWatch = append(filesToWatch, vendorFiles...)
		}
	}

	return filesToWatch, nil
}

func getFilesFromPath(config structs.Config, glob string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(config.SymfonyProjectDir, glob))
	if err != nil {
		return nil, fmt.Errorf("error while globbing files: %v", err)
	}
	return files, nil
}
