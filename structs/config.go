package structs

import "time"

// Symfony default parameters for Symfony/Flex.
const (
	EnvDefault            = "dev"
	ConsolePath           = "bin/console"
	DebugDefault          = true
	ConfigDirectory       = "config"
	TranslationsDirectory = "translations"
	TemplatesDir          = "templates"
	SrcDir                = "src"
	VendorDir             = "vendor"
	VendorDefault         = false
	SleepTime             = 30 * time.Millisecond // Watcher process sleep time
)

// DefaultExcludedDirs contains the directories that should be excluded by default.
var DefaultExcludedDirs = []string{".git", ".github", "node_modules"}

// Config holds all the parameters needed for the application. The YAML tags
// represent the keys in the Symfony custom config file, which will override
// these default values.
type Config struct {
	SymfonyProjectDir      string        `yaml:"project_dir"`      // The main Symfony project directory
	SymfonyConsolePath     string        `yaml:"console_path"`     // Relative path to the Symfony console
	SymfonyEnv             string        `yaml:"env"`              // APP_ENV parameter
	SymfonyDebug           bool          `yaml:"debug"`            // APP_DEBUG parameter
	SymfonyConfigDir       string        `yaml:"config_dir"`       // Directory where configuration files are stored
	SymfonyTranslationsDir string        `yaml:"translations_dir"` // Directory where translation files are stored
	SymfonyTemplatesDir    string        `yaml:"templates_dir"`    // Directory where template files are stored
	SymfonySrcDir          string        `yaml:"src_dir"`          // Directory where source code is stored
	SymfonyVendorDir       string        `yaml:"vendor_dir"`       // Directory where vendor code is stored
	SleepTime              time.Duration `yaml:"sleep_time"`       // Sleep time between filesystem checks
	VendorWatch            bool          `yaml:"vendor_watch"`     // Whether to watch vendor directories
	VendorList             []string      `yaml:"vendor_list"`      // List of specific vendor directories to watch
	ExcludeDirs            []string      `yaml:"exclude_dirs"`     // Directories to exclude from monitoring
}

// Init initializes the Config object with default values.
func (obj *Config) Init() {
	obj.SymfonyConsolePath = ConsolePath
	obj.SymfonyEnv = EnvDefault
	obj.SymfonyDebug = DebugDefault
	obj.SymfonyConfigDir = ConfigDirectory
	obj.SymfonyTranslationsDir = TranslationsDirectory
	obj.SymfonyTemplatesDir = TemplatesDir
	obj.SymfonySrcDir = SrcDir
	obj.SymfonyVendorDir = VendorDir
	obj.SleepTime = SleepTime
	obj.VendorWatch = VendorDefault
	obj.VendorList = []string{}
	obj.ExcludeDirs = DefaultExcludedDirs
}
