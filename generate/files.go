package generate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/sirupsen/logrus"
)

func listFiles(fileList *[]string) {
	conf := conf.Wiki

	err := os.Chdir(conf.Dir)
	if err != nil {
		logrus.WithField("module", "generate").WithError(err).Fatalf("change dir to %s failed", conf.Dir)
		return
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, exclude := range conf.Exclude {
			if strings.HasPrefix(exclude, "/") {
				// For Absolute path
				if strings.HasPrefix(path, exclude[1:]) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			} else {
				// for Name
				if info.Name() == exclude {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}

		if !info.IsDir() && containsFormat(info.Name(), conf.Format) {
			*fileList = append(*fileList, path)
		}
		return nil
	})

	if err != nil {
		logrus.WithField("module", "generate").WithError(err).Fatalf("walk dir failed")
	}
}

func containsFormat(name string, formats []string) bool {
	ext := filepath.Ext(name)
	for _, format := range formats {
		if strings.TrimPrefix(ext, ".") == format {
			return true
		}
	}
	return false
}
