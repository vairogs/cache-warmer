package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/vairogs/cache-warmer/structs"
	"github.com/vairogs/cache-warmer/symfony"
	"github.com/vairogs/cache-warmer/tools"
)

const acronym = "CacheWarmer"
const binary = "vcw"
const version = "0.0.1"
const separator = "——————————————————————————————————————————————————————————————————————"
const repository = "https://github.com/vairogs/cache-warmer"

func help() {
	fmt.Println(fmt.Sprintf("Call %s with the path of your Symfony project as the first argument.", color.New(color.FgGreen).Sprintf("%s", binary)))
	fmt.Println(fmt.Sprintf("Example: \"%s %s\"", color.New(color.FgGreen).Sprintf("%s", binary), color.New(color.FgHiYellow).Sprintf("../vairogs.com")))
	fmt.Println(fmt.Sprintf("Or even: \"%s %s\" if you call it from the root of your Symfony project directory.", color.New(color.FgGreen).Sprintf("%s", binary), color.New(color.FgHiYellow).Sprintf(".")))
	fmt.Println(fmt.Sprintf("%s %s", color.New(color.FgHiYellow).Sprintf("[💡]"), color.New(color.FgWhite).Sprintf("Add it to your $PATH if not done already.")))
}

func welcome() {
	fmt.Println(separator)
	fmt.Println(fmt.Sprintf("  %s version %s by %s - https://www.vairogs.com", color.New(color.Bold, color.FgGreen).Sprintf(acronym), color.New(color.FgHiYellow).Sprintf("v%s", version), color.New(color.FgHiRed).Sprintf("k0d3r1s")))
	fmt.Println(separator)
	fmt.Println(fmt.Sprintf("%s watches your files and automatically refreshes your project cache.", color.New(color.FgGreen).Sprintf(acronym)))
	fmt.Println(separator)
}

func errorNothingToWatch() {
	tools.PrintError(fmt.Errorf("no file to watch found"))
	fmt.Println(fmt.Sprintf("%s If you are using an \"old\" Symfony project directory structure", color.New(color.FgHiYellow).Sprintf("[💡]")))
	fmt.Println(fmt.Sprintf("     you have to customize the watched directories with a %s file", color.New(color.FgHiYellow).Sprintf(".cw.yaml")))
	fmt.Println(fmt.Sprintf("     at the root of your Symfony project. Check out the doc: %s", color.New(color.FgMagenta).Sprintf("%s", repository)))
	os.Exit(0)
}

func mainLoop(config structs.Config, filesToWatch map[string]string) {
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
		} else {
			time.Sleep(config.SleepTime)
		}
	}
}

func main() {
	var config structs.Config
	var err error
	config.Init()

	welcome()

	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}

	vendors := flag.String("vendors", "", "comma-separated list of vendors to watch")
	exclude := flag.String("exclude", "", "comma-separated directories not to watch")
	flag.Parse()

	var vendorList []string
	if *vendors != "" {
		vendorList = strings.Split(*vendors, ",")
	}

	if len(vendorList) > 0 {
		config.VendorWatch = true
		config.VendorList = vendorList
	}

	var excludeDirs []string
	if *exclude != "" {
		excludeDirs = strings.Split(*exclude, ",")
	}

	config.ExcludeDirs = append(config.ExcludeDirs, excludeDirs...)

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

	fmt.Println(" > Symfony env: " + color.New(color.FgGreen).Sprintf(strings.Trim(fmt.Sprintf("%s", out), "\n")))

	start := time.Now()
	filesToWatch, _ := symfony.GetWatchMap(config)
	end := time.Now()
	elapsed := end.Sub(start)

	if len(filesToWatch) == 0 {
		errorNothingToWatch()
	}

	fmt.Println(fmt.Sprintf(" > %s file(s) watched in %s in %s millisecond(s).", color.YellowString("%d", len(filesToWatch)), color.YellowString("%s", config.SymfonyProjectDir), color.YellowString("%d", elapsed.Milliseconds())))
	fmt.Println(fmt.Sprintf(" > %s to stop watching or run %s %s.", color.GreenString("CTRL+C"), color.GreenString("kill -9"), color.GreenString("%d", os.Getpid())))

	mainLoop(config, filesToWatch)
}
