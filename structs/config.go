package structs

import "time"

// Symfony default parameters for Symfony/Flex.
const (
	ClearCache            = false
	ConfigDirectory       = "config"
	ConsolePath           = "bin/console"
	DebugDefault          = true
	EnvDefault            = "dev"
	ForceClearCache       = false
	MigrationsDir         = "migrations"
	PoolsProvided         = false
	SleepTime             = 30 * time.Millisecond // Watcher process sleep time
	SrcDir                = "src"
	TemplatesDir          = "templates"
	TranslationsDirectory = "translations"
	VendorDefault         = false
	VendorDir             = "vendor"
)

// DefaultExcludedDirs contains the directories that should be excluded by default.
var DefaultExcludedDirs = []string{".git", ".github", "node_modules"}

// Config holds all the parameters needed for the application. The YAML tags
// represent the keys in the Symfony custom config file, which will override
// these default values.
type Config struct {
	ClearCache             bool     // Clear cache instead of only warmup
	ExcludeDirs            []string // Directories to exclude from monitoring
	ForceClearCache        bool     // Force cache removal using rm -rf var/cache
	MigrationsDir          string
	Pools                  []string      // List of pools to watch
	PoolsProvided          bool          // Whether the --pools flag was provided
	SleepTime              time.Duration // Sleep time between filesystem checks
	SymfonyConfigDir       string        // Directory where configuration files are stored
	SymfonyConsolePath     string        // Relative path to the Symfony console
	SymfonyDebug           bool          // APP_DEBUG parameter
	SymfonyEnv             string        // APP_ENV parameter
	SymfonyProjectDir      string        // The main Symfony project directory
	SymfonySrcDir          string        // Directory where source code is stored
	SymfonyTemplatesDir    string        // Directory where template files are stored
	SymfonyTranslationsDir string        // Directory where translation files are stored
	SymfonyVendorDir       string        // Directory where vendor code is stored
	VendorList             []string      // List of specific vendor directories to watch
	VendorWatch            bool          // Whether to watch vendor directories
}

// Init initializes the Config object with default values.
func (obj *Config) Init() {
	obj.ClearCache = ClearCache
	obj.ExcludeDirs = DefaultExcludedDirs
	obj.ForceClearCache = ForceClearCache
	obj.MigrationsDir = MigrationsDir
	obj.Pools = []string{}
	obj.PoolsProvided = PoolsProvided
	obj.SleepTime = SleepTime
	obj.SymfonyConfigDir = ConfigDirectory
	obj.SymfonyConsolePath = ConsolePath
	obj.SymfonyDebug = DebugDefault
	obj.SymfonyEnv = EnvDefault
	obj.SymfonySrcDir = SrcDir
	obj.SymfonyTemplatesDir = TemplatesDir
	obj.SymfonyTranslationsDir = TranslationsDirectory
	obj.SymfonyVendorDir = VendorDir
	obj.VendorList = []string{}
	obj.VendorWatch = VendorDefault
}
