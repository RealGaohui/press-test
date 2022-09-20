package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	cfg "press-test/config"
	"strconv"
	"strings"
	"time"
)

type WRK struct {
	Total       string   `json:"total"`
	Avg         string   `json:"avg"`
	Stdev       string   `json:"stdev"`
	Max         string   `json:"max"`
	P50         string   `json:"p50"`
	P75         string   `json:"p75"`
	P90         string   `json:"p90"`
	P95         string   `json:"p95"`
	P99         string   `json:"p99"`
	Connections int      `json:"connections"`
	QPS         string   `json:"qps"`
	Timeout     string   `json:"timeout"`
	Gen         []string `json:"gen"`
}

type Reporter interface {
	Report(result *Result) error
	WriteToCSV(title string) error
}

func (w *WRK) Report(result *Result) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	content := []string{"压测报告", fmt.Sprintf("时间: %s", currentTime), "配置: "}
	data, byteError := json.Marshal(*result)
	if byteError != nil {
		return byteError
	}
	tmpMap := make(map[string]interface{})
	if err = json.Unmarshal(data, &tmpMap); err != nil {
		return err
	}
	for key, value := range tmpMap {
		if key == "wrk" {
			val := value.(wrk)
			content = append(content, fmt.Sprintf(">wrk: <font color=\"comment\"> threads: %s, connections: %d} </font>", val.Threads, val.Connections))
		} else if key == "backfill" {
			val := value.(time.Duration)
			content = append(content, fmt.Sprintf("backfill: <font color=\"comment\"> 耗时%vs </font>", val.Seconds()))
		} else {
			val := value.(Resource)
			content = append(content, fmt.Sprintf(">%s: <font color=\"comment\"> replicas: %s, cpu: %s, memory: %s} </font>",
				key,
				val.Replicas,
				val.CPU,
				val.Memory,
			),
			)
		}
	}
	content = append(content, "")
	content = append(content, "压测结果: ")
	content = append(content, fmt.Sprintf(">Total Requests: <font color=\"info\"> %s </font>", w.Total))
	if w.Timeout == "" {
		w.Timeout = "0"
	}
	content = append(content, fmt.Sprintf(">Timeout: <font color=\"info\"> %s </font> ", w.Timeout))
	content = append(content, fmt.Sprintf(">QPS: <font color=\"info\"> %s </font>", w.QPS))
	content = append(content, fmt.Sprintf(">Avg: <font color=\"info\"> %s </font>", w.Avg))
	content = append(content, fmt.Sprintf(">Stdev: <font color=\"info\"> %s </font>", w.Stdev))
	content = append(content, fmt.Sprintf(">Max: <font color=\"info\"> %s </font>", w.Max))
	content = append(content, fmt.Sprintf("P50: <font color=\"info\"> %s </font>", w.P50))
	content = append(content, fmt.Sprintf("P75: <font color=\"info\"> %s </font>", w.P75))
	content = append(content, fmt.Sprintf("P90: <font color=\"info\"> %s </font>", w.P90))
	content = append(content, fmt.Sprintf("P95: <font color=\"info\"> %s </font>", w.P95))
	content = append(content, fmt.Sprintf("P99: <font color=\"info\"> %s </font>", w.P99))
	w.Gen = content
	if err = GenericFactory().SendWechat("Report", content); err != nil {
		return err
	}
	return nil
}

func (w *WRK) WriteToCSV(title string, result *Result) error {
	file := cfg.CsvFilePath + "/" + cfg.CsvFileName
	f, openFileError := os.Open(file)
	if openFileError != nil {
		return openFileError
	}
	writer := bufio.NewWriter(f)
	tmpRow := make([]string, 10)
	tmpRow = append(tmpRow, title)
	tmpRow = append(tmpRow, strconv.Itoa(result.WRK.Connections))
	tmpRow = append(tmpRow, w.QPS)
	tmpRow = append(tmpRow, w.P50)
	tmpRow = append(tmpRow, w.P75)
	tmpRow = append(tmpRow, w.P90)
	tmpRow = append(tmpRow, w.P95)
	tmpRow = append(tmpRow, w.P99)
	tmpRow = append(tmpRow, w.Avg)
	tmpRow = append(tmpRow, w.Stdev)
	tmpRow = append(tmpRow, w.Max)
	tmpRow = append(tmpRow, result.Backfill.String())
	row := strings.Join(tmpRow, ",")
	_, err = writer.WriteString(row + "\n")
	if err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	_ = f.Close()
	return nil
}
