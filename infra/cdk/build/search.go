package build

import (
	"os"
	"path/filepath"
	"regexp"
)

type searchOptions struct {
	Omit []string
}

type SearchOption func(*searchOptions)

func WithOmit(omit ...string) SearchOption {
	return func(o *searchOptions) {
		if o.Omit == nil {
			o.Omit = make([]string, 0)
		}

		o.Omit = append(o.Omit, omit...)
	}
}

// Search searches for a file in a given path and returns a list of files that match the given file name.
func Search(rootPath, fileName string, opts ...SearchOption) ([]string, error) {
	options := searchOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return search(rootPath, fileName, make([]string, 0), options)
}

func search(rootPath, fileName string, files []string, options searchOptions) ([]string, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if shouldBeOmited(entry, options) {
			continue
		}

		if entry.IsDir() {
			files, err = search(filepath.Join(rootPath, entry.Name()), fileName, files, options)
			if err != nil {
				return files, err
			}

			continue
		}

		if entry.Name() == fileName {
			files = append(files, filepath.Join(rootPath, entry.Name()))
		}
	}

	return files, nil
}

func shouldBeOmited(file os.DirEntry, opts searchOptions) bool {
	for _, omit := range opts.Omit {
		assert, err := regexp.MatchString(omit, file.Name())
		if err != nil {
			panic(err)
		}

		if assert {
			return true
		}
	}

	return false
}
