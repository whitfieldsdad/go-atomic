package atomic

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

const (
	ChannelBufferSize = 1000
)

type Agent struct {
	Id                          string `json:"id"`
	workdir                     string
	taskInvocationChannel       chan TaskInvocation
	taskInvocationStatusChannel chan TaskInvocationStatus
	artifactChannel             chan Artifact
	errorChannel                chan error
}

func NewAgent(workdir string) Agent {
	return Agent{
		Id:                          GetAgentId(),
		workdir:                     workdir,
		errorChannel:                make(chan error, ChannelBufferSize),
		artifactChannel:             make(chan Artifact, ChannelBufferSize),
		taskInvocationChannel:       make(chan TaskInvocation, ChannelBufferSize),
		taskInvocationStatusChannel: make(chan TaskInvocationStatus, ChannelBufferSize),
	}
}

func (a Agent) createTaskInvocation(e TaskInvocation) error {
	path := a.getTaskInvocationPath(e.TaskId, e.Id)
	return WriteJSONFile(path, e)
}

func (a Agent) createTaskInvocationStatus(e TaskInvocationStatus) error {
	path := a.getTaskInvocationStatusPath(e.TaskId, e.TaskInvocationId, e.StatusType.String())
	err := WriteJSONFile(path, e)
	if err != nil {
		return err
	}
	linkPath := a.getTaskInvocationStatusLinkPath(e.TaskId, e.TaskInvocationId, e.StatusType.String())
	return os.Symlink(path, linkPath)
}

func (a Agent) createTaskInvocationResult(e TaskInvocationResult) error {
	path := a.getTaskInvocationResultPath(e.TaskId, e.TaskInvocationId)
	return WriteJSONFile(path, e)
}

// TODO
func (a Agent) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go a.run(ctx, cancel, &wg)

	// Handle signals.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	log.Info("Waiting for context cancellation or SIGINT...")
	select {
	case <-ctx.Done():
		cancel()

	case <-signalChannel:
		cancel()

	}
	log.Info("Waiting for agent to stop...")
	wg.Wait()
	return nil
}

func (a Agent) run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return nil
	case taskInvocation := <-a.taskInvocationChannel:
		a.onTaskInvocation(taskInvocation)
	case taskInvocationStatus := <-a.taskInvocationStatusChannel:
		a.onTaskInvocationStatus(taskInvocationStatus)
	case artifact := <-a.artifactChannel:
		a.onArtifact(artifact)
	case err := <-a.errorChannel:
		a.onError(err)
	}
	return nil
}

func (a Agent) RunTask(taskId string) {
	invocation := NewTaskInvocation(taskId)
	a.taskInvocationChannel <- *invocation
}

func (a Agent) onTaskInvocation(taskInvocation TaskInvocation) {
	invocation := taskInvocation
	log.Infof("Running task: %s", invocation.TaskId)
}

func (a Agent) onTaskInvocationStatus(taskInvocationStatus TaskInvocationStatus) {
	//
}

func (a Agent) onArtifact(artifact Artifact) {
	// ...
}

func (a Agent) onError(err error) {
	//
}

func (a Agent) getTaskDir(taskId string) string {
	return filepath.Join(a.workdir, "tasks", taskId)
}

func (a Agent) getTaskInvocationDir(taskId, taskInvocationId string) string {
	return filepath.Join(a.workdir, "tasks", taskId, "invocations", taskInvocationId)
}

func (a Agent) getTaskInvocationStatusDir(taskId, taskInvocationId string) string {
	return filepath.Join(a.getTaskInvocationDir(taskId, taskInvocationId), "statuses")
}

func (a Agent) getTaskPath(taskId string) string {
	return filepath.Join(a.getTaskDir(taskId), "task.json")
}

func (a Agent) getTaskInvocationPath(taskId, taskInvocationId string) string {
	return filepath.Join(a.getTaskInvocationDir(taskId, taskInvocationId), "task_invocation.json")
}

func (a Agent) getTaskInvocationResultPath(taskId, taskInvocationId string) string {
	return filepath.Join(a.getTaskInvocationDir(taskId, taskInvocationId), "task_invocation_result.json")
}

func (a Agent) getTaskInvocationStatusPath(taskId, taskInvocationId, statusType string) string {
	return filepath.Join(a.getTaskInvocationStatusDir(taskId, taskInvocationId), statusType+".json")
}

func (a Agent) getTaskInvocationStatusLinkPath(taskId, taskInvocationId, statusType string) string {
	return filepath.Join(a.getTaskInvocationDir(taskId, taskInvocationId), "task_invocation_status.json")
}

func (a Agent) getArtifactDir(taskId, taskInvocationId, artifactType string) string {
	t, ok := ArtifactTypeToPlural[ArtifactType(artifactType)]
	if !ok {
		panic(fmt.Sprintf("unknown artifact type: %s", artifactType))
	}
	return filepath.Join(a.getTaskInvocationDir(taskId, taskInvocationId), "artifacts", t.String())
}

func (a Agent) getArtifactPath(taskId, taskInvocationId, artifactType, artifactId string) string {
	return filepath.Join(a.getArtifactDir(taskId, taskInvocationId, artifactType), artifactId+".json")
}

func GetAgentId() string {
	hostId := GetHostId()
	userId, err := GetCurrentUser()
	if err != nil {
		log.Fatalf("Failed to lookup current user: %v", err)
	}
	m := map[string]interface{}{
		"host_id": hostId,
		"user_id": userId,
	}
	id, err := NewUUID5FromMap(appId, m)
	if err != nil {
		log.Fatalf("Failed to generate agent ID: %v", err)
	}
	return id
}
