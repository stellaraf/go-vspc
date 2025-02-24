package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"go.stellar.af/go-utils/repo"
)

const OPENAPI_SPEC_URL string = "https://cdn.veeam.com/content/dam/helpcenter/global/reference/vspc_rest_81.yaml"

func getSpecPath(base string) (string, error) {
	specPath := filepath.Join(base, "spec.oapi3.json")
	specPath, err := filepath.Abs(specPath)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(specPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(specPath)
			if err != nil {
				return "", err
			}
			log.Printf("created %s", specPath)
			defer f.Close()
			return specPath, nil
		}
		return "", err
	}
	return specPath, nil
}

func downloadSpec(specPath string) error {
	res, err := http.Get(OPENAPI_SPEC_URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("failed to retrieve API spec: %d %s", res.StatusCode, res.Status)
		return err
	}
	specFile, err := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return err
	}
	defer specFile.Close()
	specData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	asJSON, err := yaml.YAMLToJSON(specData)
	if err != nil {
		return err
	}

	_, err = specFile.Write(asJSON)
	if err != nil {
		return err
	}
	log.Printf("wrote spec to %s", specFile.Name())
	return nil
}

func main() {
	root, err := repo.Root(4)
	if err != nil {
		panic(err)
	}
	specPath, err := getSpecPath(root)
	if err != nil {
		panic(err)
	}
	err = downloadSpec(specPath)
	if err != nil {
		panic(err)
	}
}
