package atomic

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

type Client struct{}

func NewClient() Client {
	return Client{}
}

func (c Client) ImportTaskTemplates(inputPath, outputPath string) error {
	templates, err := c.readTaskTemplatesFromPath(inputPath)
	if err != nil {
		log.Fatalf("Failed to read task templates: %v", err)
	}
	for _, template := range templates {
		filename := fmt.Sprintf("%s.json", template.Id)
		path := fmt.Sprintf("%s/%s", outputPath, filename)
		err := WriteJSONFile(path, template)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) readTaskTemplatesFromPath(inputPath string) ([]TaskTemplate, error) {
	inputPath = getRealPath(inputPath)
	st, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return c.generateTaskTemplatesFromDir(inputPath)
	} else if strings.HasSuffix(inputPath, ".tar.gz") {
		return c.generateTaskTemplatesFromTarball(inputPath)
	} else if strings.HasSuffix(inputPath, ".yaml") {
		bundle, err := readAtomicRedTeamYAMLFile(inputPath)
		if err != nil {
			return nil, err
		}
		return bundle.GetTaskTemplates()
	}
	return nil, errors.New("unsupported source")
}

func (c Client) generateTaskTemplatesFromVFS(afs afero.Afero, root string) ([]TaskTemplate, error) {
	var paths []string
	afs.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		g := glob.MustCompile("*/T*/T*.yaml")
		if g.Match(path) {
			paths = append(paths, path)
		}
		return nil
	})
	var allTemplates []TaskTemplate
	for _, path := range paths {
		f, err := afs.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		b, err := afero.ReadAll(f)
		if err != nil {
			return nil, err
		}
		bundle, err := parseAtomicRedTeamYAML(b)
		if err != nil {
			log.Warnf("Failed to parse bundle: %s - %s", path, err)
		}
		templates, err := bundle.GetTaskTemplates()
		if err != nil {
			log.Warnf("Failed to get task templates: %s - %s", path, err)
		}
		allTemplates = append(allTemplates, templates...)
	}
	return allTemplates, nil
}

func (c Client) generateTaskTemplatesFromDir(path string) ([]TaskTemplate, error) {
	afs := &afero.Afero{Fs: afero.NewOsFs()}
	return c.generateTaskTemplatesFromVFS(*afs, path)
}

func (c Client) generateTaskTemplatesFromTarball(path string) ([]TaskTemplate, error) {
	tf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tf.Close()

	gr, err := gzip.NewReader(tf)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tfs := tarfs.New(tar.NewReader(gr))
	afs := &afero.Afero{Fs: tfs}
	return c.generateTaskTemplatesFromVFS(*afs, "/")
}
