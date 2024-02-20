package atomic

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

type Client struct {
	Config ClientConfig
}

type ClientConfig struct {
	Workdir string `json:"workdir"`
}

func NewClient(workdir string) (*Client, error) {
	client := &Client{
		Config: NewClientConfig(workdir),
	}
	return client, nil
}

func NewClientConfig(workdir string) ClientConfig {
	return ClientConfig{
		Workdir: workdir,
	}
}

func (c Client) GenerateTaskTemplatesFromPath(path string) ([]TaskTemplate, error) {
	path = getRealPath(path)
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return c.generateTaskTemplatesFromDir(path)
	} else if strings.HasSuffix(path, ".tar.gz") {
		return c.generateTaskTemplatesFromTarball(path)
	} else if strings.HasSuffix(path, ".yaml") {
		bundle, err := readAtomicRedTeamYAMLFile(path)
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
