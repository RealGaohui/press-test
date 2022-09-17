package utils

import (
	"os"
	"path/filepath"
	cfg "press-test/config"
)

func CreateFile() error {
	err = initPath(cfg.CsvFilePath)
	if err != nil {
		return err
	}
	err = initPath(cfg.WrkRawlogPath)
	if err != nil {
		return err
	}
	err = initPath(cfg.SchedulerLogPath)
	if err != nil {
		return err
	}
	return nil
}

func initPath(path string) error {
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Log.Errorf("Failed to create directory %s", path)
			return err
		}
	}
	return nil
}

func initFile(path string) error {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = create(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func create(path string) error {
	s2 := filepath.Join(path, "../")
	s3, err1 := filepath.Abs(s2)
	if err1 != nil {
		return err1
	}
	err = os.MkdirAll(s3, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.Create(path)
	if err != nil {
		return err
	}
	return nil
}
