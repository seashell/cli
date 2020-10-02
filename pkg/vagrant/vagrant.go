package vagrant

import (
	"fmt"
	"os"
	"os/exec"
)

type Vagrant struct {
	path       string
	executable string
}

func New(path string) (*Vagrant, error) {
	e, err := exec.LookPath("vagrant")
	if err != nil {
		return nil, err
	}

	return &Vagrant{
		path:       path,
		executable: string(e),
	}, nil
}

func (v *Vagrant) Up() error {
	cmd := exec.Command(v.executable, "up")
	cmd.Dir = v.path
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (v *Vagrant) Destroy() error {
	cmd := exec.Command(v.executable, "destroy", "-f")
	cmd.Dir = v.path
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (v *Vagrant) Reload() error {
	cmd := exec.Command(v.executable, "reload")
	cmd.Dir = v.path
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (v *Vagrant) Provision() error {
	cmd := exec.Command(v.executable, "provision")
	cmd.Dir = v.path
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (v *Vagrant) SSH(envVars map[string]string) error {
	env := ""
	if envVars != nil {
		env = "export"
		for key, element := range envVars {
			env = fmt.Sprintf("%s %s=%s", env, key, element)
		}
		env = fmt.Sprintf("%s;", env)
	}

	command := fmt.Sprintf("cd /workspace && %s sudo -E su", env)
	cmd := exec.Command(v.executable, "ssh", "-c", command)
	cmd.Dir = v.path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
