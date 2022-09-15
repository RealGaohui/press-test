package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	cfg "press-test/config"
	"sync"
)

var (
	f     *os.File
	Error error
	err   error
)

type Log struct {
	lock sync.Mutex
	file *string
}

type Interfaceas interface {
	Info(args interface{})
	Infof(format string, args interface{})
	Warn(args interface{})
	Warnf(format string, args interface{})
	Error(args interface{})
	Errorf(format string, args interface{})
	InitLogFile() error
}

func Logger() *Log {
	file := cfg.LogFile
	return &Log{file: &file}
}

func (l *Log) Info(args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintln(args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer func() {
		_ = f.Close()
	}()
}

func (l *Log) Infof(format string, args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf(format, args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer f.Close()
}

func (l *Log) Warn(args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintln(args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer f.Close()
}

func (l *Log) Warnf(format string, args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf(format, args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer f.Close()
}

func (l *Log) Error(args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintln(args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer f.Close()
}

func (l *Log) Errorf(format string, args interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f, err = os.OpenFile(*l.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf(format, args)
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(msg + "\n")
	if err != nil {
		panic(err)
	}
	_ = writer.Flush()
	defer f.Close()
}

func (l *Log) InitLogFile() error {
	_, err := os.Stat(*l.file)
	if os.IsNotExist(err) {
		err = createFile(*l.file)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(path string) error {
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
