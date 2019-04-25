package main

import (
	"github.com/openland/baikonur/il"
	"os"
	"path/filepath"
	"strings"
)

func collectFiles(path string) ([]string, error) {
	var res []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".graphql") {
			return nil
		}
		res = append(res, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	files, err := collectFiles("tests/queries")
	if err != nil {
		panic(err)
	}
	model := il.LoadModel("tests/schema.json", files)
	print(len(model.Fragments))
	// generate(model)
}
