package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/vairogs/cache-warmer/structs"
	"github.com/vairogs/cache-warmer/symfony"
	"github.com/vairogs/cache-warmer/tools"
)

var version = "nightly"

const (
	acronym = "CacheWarmer"
	//binary     = "vcw"
	repository = "https://github.com/vairogs/cache-warmer"
)

// Help prints the instructions on how to use the `VCW` command line tool. It provides examples of how to call the command, including the required argument for the path of the Symfony project. It also suggests adding the command to the system's `$PATH` if it has not been done already.
func Help() {
	binary := filepath.Base(os.Args[0])
	binary = strings.TrimSuffix(binary, filepath.Ext(binary))

	fmt.Println(fmt.Sprintf("Call %s with the path of your Symfony project as the first argument.", color.New(color.FgGreen).Sprintf("%s", binary)))
	fmt.Println(fmt.Sprintf("Example: \"%s %s\"", color.New(color.FgGreen).Sprintf("%s", binary), color.New(color.FgHiYellow).Sprintf(".")))
	fmt.Println(fmt.Sprintf("Or even: \"bin/%s %s\" if you call it from the bin of your Symfony project directory.", color.New(color.FgGreen).Sprintf("%s", binary), color.New(color.FgHiYellow).Sprintf(".")))
}

// Welcome prints a welcome message to the console.
// It retrieves the version and generates a clickable version link.
// The message includes the acronym, version, and author's website.
// It also provides a brief description of the functionality of the program.
func Welcome() {
	versionURL := GenerateVersionLink(version)
	clickableVersion := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", versionURL, version)

	var length = 80
	fmt.Println(GenerateSeparator(length))
	fmt.Println(fmt.Sprintf("  %s version %s by %s - https://me.k0d3r1s.com", color.New(color.Bold, color.FgGreen).Sprintf(acronym), color.New(color.FgHiYellow).Sprintf(clickableVersion), color.New(color.FgHiRed).Sprintf("k0d3r1s")))
	fmt.Println(GenerateSeparator(length))
	fmt.Println(fmt.Sprintf("%s watches your files and automatically refreshes your project cache.", color.New(color.FgGreen).Sprintf(acronym)))
	fmt.Println(GenerateSeparator(length))
}

// ErrorNothingToWatch displays an error message indicating that no files to watch were found.
func ErrorNothingToWatch() {
	tools.PrintError(fmt.Errorf("no file to watch found"))
	fmt.Println(fmt.Sprintf("%s If you are using an \"old\" Symfony project directory structure", color.New(color.FgHiYellow).Sprintf("[💡]")))
	fmt.Println(fmt.Sprintf("     you have to customize the watched directories with a %s file", color.New(color.FgHiYellow).Sprintf(".cw.yaml")))
	fmt.Println(fmt.Sprintf("     at the root of your Symfony project. Check out the doc: %s", color.New(color.FgMagenta).Sprintf("%s", repository)))
	os.Exit(0)
}

// MainLoop continuously monitors for file changes and performs cache warming if an update is detected.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values for the application.
// It also takes a `filesToWatch` parameter of type map[string]string that represents the files to watch for changes.
// The function checks for updated files using `symfony.GetWatchMap` and compares it with the existing `filesToWatch` map.
// If there are any differences, it starts cache warming by calling `symfony.CacheWarmup`. It measures the time taken
// to warm up the cache and prints the result. The updated `filesToWatch` map is then assigned to `filesToWatch`.
// If there are no differences, the function sleeps for a specified duration defined in the `config` parameter.
func MainLoop(config structs.Config, filesToWatch map[string]string) {
	for {
		updatedFiles, _ := symfony.GetWatchMap(config)
		if !reflect.DeepEqual(filesToWatch, updatedFiles) {
			start := time.Now()
			fmt.Println(fmt.Sprintf(" %s at %s > refreshing cache...", color.New(color.FgHiYellow).Sprintf("⬇ Update detected"), color.New(color.FgGreen).Sprintf(start.Format("15:04:05"))))
			_, _ = symfony.CacheWarmup(config)
			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Println(fmt.Sprintf("  %s in %s second(s).", color.New(color.FgGreen).Sprintf("✅  Done!"), color.New(color.FgHiYellow).Sprintf("%.2f", elapsed.Seconds())))
			filesToWatch = updatedFiles
			fmt.Println(fmt.Sprintf(" > %s file(s) watched in %s", color.YellowString("%d", len(filesToWatch)), color.YellowString("%s", config.SymfonyProjectDir)))
		} else {
			time.Sleep(config.SleepTime)
		}
	}
}

// main is the entry point of the program. It initializes the configuration, displays a Welcome message,
// parses command line arguments, checks for required parameters, sets configuration values based on the command line arguments,
// checks for the existence of Symfony console, retrieves the Symfony version, gets the files to watch,
// displays some information about the project, and enters the main loop to monitor and react to file changes.
func main() {
	var config structs.Config
	var err error
	config.Init()

	Welcome()

	if len(os.Args) == 1 {
		Help()
		os.Exit(0)
	}

	vendors := flag.String("vendor", "", "comma-separated list of vendors to watch")
	exclude := flag.String("exclude", "", "comma-separated directories not to watch")
	clearCache := flag.Bool("cache", false, "clear cache instead of just warmup")
	forceClearCache := flag.Bool("force", false, "force clear cache (rm -rf var/cache)")
	noDebug := flag.Bool("no-debug", false, "force clear cache (rm -rf var/cache)")
	env := flag.String("env", "dev", "comma-separated list of vendors to watch")

	pools := structs.NewCustomFlag()
	flag.Var(pools, "pools", "comma-separated list of pools to clear")

	flag.Parse()

	config.SymfonyEnv = *env
	config.ClearCache = *clearCache
	config.SymfonyDebug = !*noDebug

	if *forceClearCache {
		config.ClearCache = false
		config.ForceClearCache = true
	}

	if *vendors != "" {
		vendorList := ParseCommaSeparated(*vendors)
		if len(vendorList) > 0 {
			config.VendorWatch = true
			config.VendorList = vendorList
		}
	}

	if *exclude != "" {
		excludeDirs := ParseCommaSeparated(*exclude)
		config.ExcludeDirs = append(config.ExcludeDirs, excludeDirs...)
	}

	if pools.IsChanged() {
		config.PoolsProvided = true
		config.Pools = pools.Get()

		if len(config.Pools) == 0 {
			config.Pools = []string{"--all"}
		}
	}

	config.SymfonyProjectDir, err = symfony.GetSymfonyProjectDir()

	if err != nil {
		tools.PrintError(fmt.Errorf("project directory not found"))
		tools.PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Project directory: " + color.New(color.FgGreen).Sprintf(config.SymfonyProjectDir))

	err = symfony.CheckSymfonyConsole(config)
	if err != nil {
		tools.PrintError(fmt.Errorf("symfony console not found"))
		tools.PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Symfony console path: " + color.New(color.FgGreen).Sprintf(config.SymfonyConsolePath))

	out, err := symfony.Version(config)
	if err != nil {
		tools.PrintError(fmt.Errorf("error while running the Symfony version command"))
		tools.PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Symfony env: " + color.New(color.FgGreen).Sprintf(strings.TrimSpace(fmt.Sprintf("%s", out))))

	start := time.Now()
	filesToWatch, _ := symfony.GetWatchMap(config)
	end := time.Now()
	elapsed := end.Sub(start)

	if len(filesToWatch) == 0 {
		ErrorNothingToWatch()
	}

	fmt.Println(fmt.Sprintf(" > %s file(s) watched in %s in %s millisecond(s).", color.YellowString("%d", len(filesToWatch)), color.YellowString("%s", config.SymfonyProjectDir), color.YellowString("%d", elapsed.Milliseconds())))
	fmt.Println(fmt.Sprintf(" > %s to stop watching or run %s %s.", color.GreenString("CTRL+C"), color.GreenString("kill -9"), color.GreenString("%d", os.Getpid())))

	MainLoop(config, filesToWatch)
}

// ParseCommaSeparated splits a comma-separated input string and returns an array of strings.
// If the input string is empty, it returns an empty array.
func ParseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	return strings.Split(input, ",")
}

// GenerateSeparator returns a string consisting of a specified number of em dashes.
// The length parameter determines the number of em dashes in the output string.
// The function uses the strings.Repeat function to repeat the em dash character.
func GenerateSeparator(length int) string {
	return strings.Repeat("—", length)
}

// IsSemanticVersion checks if a given version string is in the semantic version format.
func IsSemanticVersion(v string) bool {
	semVerPattern := `^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$`
	matched, _ := regexp.MatchString(semVerPattern, v)

	return matched
}

// GenerateVersionLink generates a link for a given version string.
// If the version is a semantic version, it creates a link to the release tag on the repository.
// If the version is not a semantic version, it creates a link to the commit on the repository.
// It returns the generated link as a string.
func GenerateVersionLink(ver string) string {
	if IsSemanticVersion(ver) || "nightly" == ver {
		return fmt.Sprintf("%s/releases/tag/%s", repository, ver)
	}

	return fmt.Sprintf("%s/commit/%s", repository, ver)
}
