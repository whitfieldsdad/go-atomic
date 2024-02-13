package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"
	"slices"

	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
	"gopkg.in/yaml.v3"
)

func getClient(flags *pflag.FlagSet) (*atomic.Client, error) {
	workdir, _ := flags.GetString("workdir")
	client, err := atomic.NewClient(workdir)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getTaskTemplateQuery(flags *pflag.FlagSet) (*atomic.TaskTemplateQuery, error) {
	taskTemplateIds, _ := flags.GetStringArray("task-template-id")
	attackTechniqueIds, _ := flags.GetStringArray("attack-technique-id")
	platforms, _ := flags.GetStringArray("platform")
	tags, _ := flags.GetStringArray("tag")

	query := &atomic.TaskTemplateQuery{
		TaskTemplateIds:    taskTemplateIds,
		AttackTechniqueIds: attackTechniqueIds,
		Platforms:          platforms,
		Tags:               tags,
	}
	auto, _ := flags.GetBool("auto")
	if auto {
		if !slices.Contains(query.Platforms, runtime.GOOS) {
			query.Platforms = append(query.Platforms, runtime.GOOS)
		}
		elevated := atomic.IsElevated()
		if !elevated {
			elevationRequired := false
			query.ElevationRequired = &elevationRequired
		}
	}
	if flags.Changed("elevation-required") {
		elevationRequired, _ := flags.GetBool("elevation-required")
		query.ElevationRequired = &elevationRequired
	}
	return query, nil
}

func printJson(v interface{}) {
	blob, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(blob))
}

func printYaml(v interface{}) {
	blob, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(blob))
}
