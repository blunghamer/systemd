package systemd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

const userflag = "--user"
const systemunits = "/lib/systemd/system"
const userunits = "/lib/systemd/user"

func Enable(unit string, user bool) error {
	us := ""
	if user {
		us = userflag
	}
	task := execute.ExecTask{Command: "systemctl",
		Args:        []string{"enable", us, unit},
		StreamStdio: false,
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("error executing task %s %v, stderr: %s", task.Command, task.Args, res.Stderr)
	}

	return nil
}

func Start(unit string, user bool) error {
	us := ""
	if user {
		us = userflag
	}
	task := execute.ExecTask{Command: "systemctl",
		Args:        []string{"start", us, unit},
		StreamStdio: false,
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("error executing task %s %v, stderr: %s", task.Command, task.Args, res.Stderr)
	}

	return nil
}

func DaemonReload() error {
	task := execute.ExecTask{Command: "systemctl",
		Args:        []string{"daemon-reload"},
		StreamStdio: false,
	}

	res, err := task.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("error executing task %s %v, stderr: %s", task.Command, task.Args, res.Stderr)
	}

	return nil
}

func InstallUnit(name string, tplFile string, tokens map[string]string, user bool) error {
	if len(tokens["Cwd"]) == 0 {
		return fmt.Errorf("key Cwd expected in tokens parameter")
	}

	tmpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("error loading template %s, error %s", tplFile, err)
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, tokens)
	if err != nil {
		return err
	}

	err = InstallUnitFile(name, tpl.Bytes(), user)
	if err != nil {
		return err
	}
	return nil
}

func InstallUnitFile(name string, unitfile []byte, user bool) error {
	err := writeUnit(name+".service", unitfile, user)
	if err != nil {
		return err
	}
	return nil
}

func writeUnit(name string, data []byte, user bool) error {
	unitpath := systemunits
	if user {
		unitpath = userunits
	}

	f, err := os.Create(filepath.Join(unitpath, name))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
