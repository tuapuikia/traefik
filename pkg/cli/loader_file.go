package cli

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/file"
	"github.com/containous/traefik/v2/pkg/config/flag"
	"github.com/containous/traefik/v2/pkg/log"
)

// FileLoader loads a configuration from a file.
type FileLoader struct {
	ConfigFileFlag string
	filename       string
}

// GetFilename returns the configuration file if any.
func (f *FileLoader) GetFilename() string {
	return f.filename
}

// Load loads the command's configuration from a file either specified with the -traefik.configfile flag, or from default locations.
func (f *FileLoader) Load(args []string, cmd *Command) (bool, error) {
	ref, err := flag.Parse(args, cmd.Configuration)
	if err != nil {
		_ = cmd.PrintHelp(os.Stdout)
		return false, err
	}

	configFileFlag := "traefik.configfile"
	if _, ok := ref["traefik.configFile"]; ok {
		configFileFlag = "traefik.configFile"
	}

	if f.ConfigFileFlag != "" {
		configFileFlag = "traefik." + f.ConfigFileFlag
		if _, ok := ref[strings.ToLower(configFileFlag)]; ok {
			configFileFlag = "traefik." + strings.ToLower(f.ConfigFileFlag)
		}
	}

	configFile, err := loadConfigFiles(ref[configFileFlag], cmd.Configuration)
	if err != nil {
		return false, err
	}

	f.filename = configFile

	if configFile == "" {
		return false, nil
	}

	logger := log.WithoutContext()
	logger.Printf("Configuration loaded from file: %s", configFile)

	content, _ := ioutil.ReadFile(configFile)
	logger.Debug(string(content))

	return true, nil
}

// loadConfigFiles tries to decode the given configuration file and all default locations for the configuration file.
// It stops as soon as decoding one of them is successful.
func loadConfigFiles(configFile string, element interface{}) (string, error) {
	finder := Finder{
		BasePaths:  []string{"/etc/traefik/traefik", "$XDG_CONFIG_HOME/traefik", "$HOME/.config/traefik", "./traefik"},
		Extensions: []string{"toml", "yaml", "yml"},
	}

	filePath, err := finder.Find(configFile)
	if err != nil {
		return "", err
	}

	if len(filePath) == 0 {
		return "", nil
	}

	if err = file.Decode(filePath, element); err != nil {
		return "", err
	}
	return filePath, nil
}
