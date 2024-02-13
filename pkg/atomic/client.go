package atomic

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func (c ClientConfig) GetTaskTemplateDir() string {
	return filepath.Join(c.Workdir, "task_templates")
}

func (c ClientConfig) GetTaskTemplatePath(taskTemplateId string) string {
	filename := fmt.Sprintf("%s.json", taskTemplateId)
	return filepath.Join(c.GetTaskTemplateDir(), filename)
}

func (c ClientConfig) GetTaskDir() string {
	return filepath.Join(c.Workdir, "tasks")
}

func (c ClientConfig) GetTaskPath(taskId string) string {
	filename := fmt.Sprintf("%s.json", taskId)
	return filepath.Join(c.GetTaskDir(), filename)
}

func (c ClientConfig) GetTaskInvocationDir() string {
	return filepath.Join(c.Workdir, "task_invocations")
}

func (c ClientConfig) GetTaskInvocationPath(taskId string) string {
	filename := fmt.Sprintf("%s.json", taskId)
	return filepath.Join(c.GetTaskInvocationDir(), filename)
}

func (c ClientConfig) GetTaskInvocationStatusDir() string {
	return filepath.Join(c.Workdir, "task_invocation_statuses")
}

func (c ClientConfig) GetTaskInvocationStatusPath(taskInvocationId string) string {
	filename := fmt.Sprintf("%s.json", taskInvocationId)
	return filepath.Join(c.GetTaskInvocationStatusDir(), filename)
}

func (c ClientConfig) GetTaskInvocationResultDir() string {
	return filepath.Join(c.Workdir, "task_invocation_results")
}

func (c ClientConfig) GetTaskInvocationResultPath(taskInvocationId string) string {
	filename := fmt.Sprintf("%s.json", taskInvocationId)
	return filepath.Join(c.GetTaskInvocationResultDir(), filename)
}

func (c ClientConfig) GetArtifactDir(artifactType string) string {
	return filepath.Join(c.Workdir, "artifacts", artifactType)
}

func (c ClientConfig) GetArtifactPath(artifactType, artifactId string) string {
	filename := fmt.Sprintf("%s.json", artifactId)
	return filepath.Join(c.GetArtifactDir(artifactType), filename)
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

func (c Client) ListTaskTemplates(q *TaskTemplateQuery) ([]TaskTemplate, error) {
	var templates []TaskTemplate
	paths, err := findFiles(c.Config.GetTaskTemplateDir(), "*.json")
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		template, err := c.readTaskTemplate(path)
		if err != nil {
			log.Errorf("Failed to read task: %s", err)
			continue
		}
		if q != nil && !q.Matches(*template) {
			continue
		}
		templates = append(templates, *template)
	}
	return templates, nil
}

func (c Client) CountTaskTemplates(q *TaskTemplateQuery) (int, error) {
	templates, err := c.ListTaskTemplates(q)
	if err != nil {
		return 0, err
	}
	return len(templates), nil
}

func (c Client) readTaskTemplate(path string) (*TaskTemplate, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var template TaskTemplate
	err = json.Unmarshal(b, &template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (c Client) ImportTaskTemplates(paths ...string) error {
	templateDir := c.Config.GetTaskTemplateDir()
	for _, root := range paths {
		templates, err := c.generateTaskTemplates(root)
		if err != nil {
			return err
		}
		log.Infof("Importing %d task templates from %s to %s", len(templates), root, templateDir)
		for _, template := range templates {
			filename := fmt.Sprintf("%s.json", template.Id)
			path := filepath.Join(templateDir, filename)
			err = writeJSONFile(path, template)
			if err != nil {
				return err
			}
		}
		log.Infof("Imported %d task templates from %s to %s", len(templates), root, templateDir)
	}
	return nil
}

func (c Client) generateTaskTemplates(path string) ([]TaskTemplate, error) {
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

func (c Client) generateTaskTemplatesFromAFS(afs afero.Afero, root string) ([]TaskTemplate, error) {
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
	return c.generateTaskTemplatesFromAFS(*afs, path)
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
	return c.generateTaskTemplatesFromAFS(*afs, "/")
}
