package dev

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"sync"
	"text/template"

	log "github.com/seashell/cli/pkg/log"
	vagrant "github.com/seashell/cli/pkg/vagrant"
)

// Dev :
type Dev struct {
	config       *Config
	logger       log.Logger
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New creates a new Dev struct from the configuration,
// potentially returning an error
func New(config *Config, logger log.Logger) (*Dev, error) {

	if logger == nil {
		return nil, errors.New("missing logger")
	}

	config.VagrantDir = path.Join(config.DataDir, config.ProjectID)

	d := &Dev{
		config:     config,
		logger:     logger,
		shutdownCh: make(chan struct{}),
	}

	return d, nil
}

// Init :
func (d *Dev) Init() error {

	//create project dir
	projectDir := path.Join(d.config.DataDir, d.config.ProjectID)
	err := os.MkdirAll(projectDir, os.ModePerm)
	if err != nil {
		return err
	}

	vagrantFileTemplate := `
Vagrant.configure("2") do |config|
	config.vm.define "{{ .ProjectID }}"
	config.vm.hostname = "{{ .ProjectID }}"
	config.vm.box = "seashell/dev"
	config.vm.box_version = "0.0.1"
	config.vm.synced_folder "{{ .ProjectDir }}", "/workspace"
end
`
	//populate vagrantfile
	t, err := template.New("vagrant").Parse(vagrantFileTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(d.config.VagrantDir, "Vagrantfile"))
	if err != nil {
		return err
	}

	err = t.Execute(f, d.config)
	if err != nil {
		return err
	}

	v, err := vagrant.New(d.config.VagrantDir)
	if err != nil {
		return err
	}

	if err := v.Up(); err != nil {
		return err
	}

	//TODO: Fetch ProjectSettings from somewhere
	// could be a config file, if API is too complex

	env := map[string]string{
		"NOMAD_ADDR": "http://35.207.108.145:4646",
	}

	if err := v.SSH(env); err != nil {
		return err
	}

	return nil
}

// Destroy :
func (d *Dev) Destroy() error {

	v, err := vagrant.New(d.config.VagrantDir)
	if err != nil {
		return err
	}

	if err := v.Destroy(); err != nil {
		return err
	}

	if err := os.RemoveAll(d.config.VagrantDir); err != nil {
		return err
	}

	return nil
}

// List :
func (d *Dev) List() error {
	cmd := exec.Command("ls", ".")
	cmd.Dir = d.config.DataDir
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Update :
func (d *Dev) Update() error {
	v, err := vagrant.New(d.config.VagrantDir)
	if err != nil {
		return err
	}

	if err := v.Reload(); err != nil {
		return err
	}

	if err := v.Destroy(); err != nil {
		return err
	}

	if err := v.Up(); err != nil {
		return err
	}

	return nil
}
