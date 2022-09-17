package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	cfg "press-test/config"
	"runtime"
)

type PV struct {
	file  *string
	index *int
}

func pv(file *string) *PV {
	return &PV{
		file: file,
	}
}

type InitPV interface {
	Shell
	ChangePVName(name, newfile string) *PV
	ChangePVPath(pvname, path, newfile string) *PV
	ChangePVHost(pvName, host, newfile string) *PV
	ChangeStorageClassName(pvname, ns string) *PV
}

func (p *PV) ChangePVConfig(hostNums, index int) *PV {
	return &PV{}
}

func (p *PV) ChangePVName(name string) *PV {
	executeCommand := ""
	if runtime.GOOS == "darwin" {
		executeCommand = fmt.Sprintf("sed -i \"\" -e 's/pv-name/%s/g' %s", name, *p.file)
	} else {
		executeCommand = fmt.Sprintf("sed -i 's/pv-name/%s/g' %s", name, *p.file)
	}
	_, err = p.Command(executeCommand)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Warnf("failed to change pv name for %s", name)
		return &PV{}
	}
	return &PV{}
}

func (p *PV) ChangePVPath(pvname, path string) *PV {
	executeCommand := ""
	if runtime.GOOS == "darwin" {
		executeCommand = fmt.Sprintf("sed -i \"\" -e 's/homePath/%s/g' %s", path, *p.file)
	} else {
		executeCommand = fmt.Sprintf("sed -i 's/homePath/%s/g' %s", path, *p.file)
	}
	_, err = p.Command(executeCommand)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Warnf("failed to change pv host path for %s", pvname)
		return &PV{}
	}
	return &PV{}
}

func (p *PV) ChangePVHost(pvname, host string) *PV {
	executeCommand := ""
	if runtime.GOOS == "darwin" {
		executeCommand = fmt.Sprintf("sed -i \"\" -e 's/server/%s/g' %s", host, *p.file)
	} else {
		executeCommand = fmt.Sprintf("sed -i 's/server/%s/g' %s", host, *p.file)
	}
	_, err = p.Command(executeCommand)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Warnf("failed to change pv host for %s", pvname)
		return &PV{}
	}
	return &PV{}
}

func (p *PV) ChangeStorageClassName(pvname, ns string) *PV {
	scn := "cassandra-" + ns
	executeCommand := ""
	if runtime.GOOS == "darwin" {
		executeCommand = fmt.Sprintf("sed -i \"\" -e 's/undefinedStorageName/%s/g' %s", scn, *p.file)
	} else {
		executeCommand = fmt.Sprintf("sed -i 's/undefinedStorageName/%s/g' %s", scn, *p.file)
	}
	_, err = p.Command(executeCommand)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Warnf("failed to change storage class name for %s", pvname)
		return &PV{}
	}
	return &PV{}
}

func (p *PV) Command(cmd string) (string, error) {
	command = exec.Command("bash", "-c", cmd)
	output, err = command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
