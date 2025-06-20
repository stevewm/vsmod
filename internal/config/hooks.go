package config

import (
	"bytes"
	"os/exec"
	"runtime"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Hook string

type CommandHooks struct {
	Pre_Run  Hook `yaml:"pre_run"`
	Post_Run Hook `yaml:"post_run"`
}

func (h Hook) Run(data any) error {
	tmpl, err := template.New("hook").Parse(string(h))
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	log.Debugf("Running hook: \"%s", buf.String()+"\"")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", buf.String())
		cmd.Stdout = log.StandardLogger().Out
		cmd.Stderr = log.StandardLogger().Out
		return cmd.Run()
	}
	cmd := exec.Command("sh", "-c", buf.String())
	cmd.Stdout = log.StandardLogger().Out
	cmd.Stderr = log.StandardLogger().Out

	return cmd.Run()
}
