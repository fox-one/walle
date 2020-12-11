package cmdutil

import (
	"os"
	"path"

	"github.com/fox-one/pkg/config"
	"github.com/mitchellh/go-homedir"
)

func HomePath(filename string) string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return path.Join(dir, filename)
}

func LoadConfig(v interface{}, filenames ...string) error {
	var filename string
	for _, name := range filenames {
		if isFileExists(name) {
			filename = name
			break
		}
	}

	config.AutomaticLoadEnv("WALLE")
	return config.LoadYaml(filename, v)
}

func isFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
