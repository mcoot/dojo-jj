package dependencies

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mcoot/dojo-jj/internal/models"
)

const (
	jjBinary               = "jj"
	jjWorkspaceRefTemplate = "name ++ \"\\t\" ++ target.change_id() ++ \"\\t\" ++ root ++ \"\\n\""
)

type JJClient interface {
	IsJJAvailable() bool
	ListWorkspaces() ([]*models.JJWorkspace, error)
	AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error)
}

type execCommandFunc func(name string, arg ...string) *exec.Cmd

type JJClientImpl struct{}

func NewJJClient() *JJClientImpl {
	return &JJClientImpl{}
}

func (c *JJClientImpl) IsJJAvailable() bool {
	_, err := exec.LookPath(jjBinary)
	return err == nil
}

func (c *JJClientImpl) ListWorkspaces() ([]*models.JJWorkspace, error) {
	output, err := c.runJJ("workspace", "list", "-T", jjWorkspaceRefTemplate)
	if err != nil {
		return nil, err
	}

	return parseJJWorkspaces(string(output))
}

func (c *JJClientImpl) AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error) {
	return nil, nil
}

func (c *JJClientImpl) runJJ(args ...string) ([]byte, error) {
	cmd := exec.Command(jjBinary, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, jjCommandError(args, stderr.String(), err)
	}

	return stdout.Bytes(), nil
}

func jjCommandError(args []string, stderr string, err error) error {
	command := jjBinary
	if len(args) > 0 {
		command += " " + strings.Join(args, " ")
	}

	stderr = strings.TrimSpace(stderr)
	if stderr == "" {
		return fmt.Errorf("%s failed: %w", command, err)
	}

	return fmt.Errorf("%s failed: %w: %s", command, err, stderr)
}

func parseJJWorkspaces(output string) ([]*models.JJWorkspace, error) {
	output = strings.TrimSuffix(output, "\n")
	if output == "" {
		return []*models.JJWorkspace{}, nil
	}

	lines := strings.Split(output, "\n")
	workspaces := make([]*models.JJWorkspace, 0, len(lines))

	for index, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			return nil, fmt.Errorf(
				"parse jj workspace list: line %d: expected 3 tab-separated fields, got %d",
				index+1,
				len(fields),
			)
		}

		workspaces = append(workspaces, &models.JJWorkspace{
			Name:     fields[0],
			ChangeID: fields[1],
			Root:     fields[2],
		})
	}

	return workspaces, nil
}
