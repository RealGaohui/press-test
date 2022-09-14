package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/RealGaohui/urlBuilder"
	"github.com/tidwall/gjson"
	"go/src/strconv"
	"io/ioutil"
	"net/http"
	cfg "press-test/config"
	logger "press-test/log"
)



type BackfillRequest struct {
	Name                 string        `json:"name"`
	VelocityAccumulators []interface{} `json:"velocityAccumulators"`
	PipelineS3Path       interface{}   `json:"pipelineS3Path"`
	DfePathList          []interface{} `json:"dfePathList"`
	ReplayTask           replayTask    `json:"replayTask"`
	ParentId             int           `json:"parentId"`
}

type replayTask struct {
	Features                   []interface{} `json:"features"`
	PipelineS3Path             interface{}   `json:"pipelineS3Path"`
	ReplayFileList             []string      `json:"replayFileList"`
	DataSourceDataSetRelations []interface{} `json:"dataSourceDataSetRelations"`
	Creator                    string        `json:"creator"`
	OutputSuffix               interface{}   `json:"outputSuffix"`
	OutputFilePathList         []interface{} `json:"outputFilePathList"`
	ReplayFileType             string        `json:"replayFileType"`
	OutputFileType             interface{}   `json:"outputFileType"`
	Delimiter                  string        `json:"delimiter"`
	EventTimeName              string        `json:"eventTimeName"`
	Encoding                   string        `json:"encoding"`
	UserIdFeatureName          interface{}   `json:"userIdFeatureName"`
	Rules                      []interface{} `json:"rules"`
	RuleSets                   []interface{} `json:"ruleSets"`
	DatasetId                  interface{}   `json:"datasetId"`
	EventTimeFormat            string        `json:"eventTimeFormat"`
	HistoricDataset            interface{}   `json:"historicDataset"`
	ClickHouseTTL              int           `json:"clickHouseTTL"`
	CodeVersionMap             struct{}      `json:"codeVersionMap"`
	IsRuleReplay               bool          `json:"isRuleReplay"`
	InsertEvent                bool          `json:"insertEvent"`
	ProcessAllAccumulator      bool          `json:"processAllAccumulator"`
	OutputEventId              bool          `json:"outputEventId"`
	InsertOnly                 bool          `json:"insertOnly"`
	SkipWriteClickhouse        bool          `json:"skipWriteClickhouse"`
	ValidateEvents             bool          `json:"validateEvents"`
	ValidateFailThreshold      float64       `json:"validateFailThreshold"`
	ReplayMode                 string        `json:"replayMode"`
	ClusterName                string        `json:"clusterName"`
}


type Backfill struct {
	backfill BackfillRequest
	taskId  int
}

type ID int

type check interface {
	CheckBackfillFinish() bool
}

type Interface interface {
	DoBackfill() (int, error)
	Alert
}

func NewID(id int)check{
	i := ID(id)
	return &i
}

func (i *ID)CheckBackfillFinish() bool {
	url := urlBuilder.URLBuilder().
		SetBase(cfg.FpBaseUrl).
		SetPath("/" ).
		SetPath(cfg.Client).
		SetPath("/backfill/").
		SetPath(strconv.Itoa(int(*i))).
		ToString()
	request, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		return newRequestErr == nil
	}
	request.Header.Set("Content-Type", "application/json")
	response, HttpDoError := http.DefaultClient.Do(request)
	if HttpDoError != nil {
		return HttpDoError == nil
	}
	resp, ReadBodyError := ioutil.ReadAll(response.Body)
	if ReadBodyError != nil {
		return ReadBodyError == nil
	}
	if !gjson.Valid(string(resp)){
		return false
	}
	if gjson.Get(string(resp), "status").String() == "COMPLETED" {
		logger.Logger().Infof("task %d compelted", int(*i))
		return true
	}
	return false
}

func (b *Backfill)DoBackfill() (check, error){
	var Check check
	logger.Logger().Infof("create backfill task: %s", b.backfill.Name)
	URL := urlBuilder.URLBuilder().
		SetBase(cfg.FpBaseUrl).
		SetPath("/").
		SetPath(cfg.Client).
		SetPath("/backfill/distribute/create_and_run").
		ToString()
	body, err := json.Marshal(b.backfill)
	if err != nil {
		return Check, err
	}
	request, newRequestErr := http.NewRequest("POST", URL, bytes.NewReader(body))
	if newRequestErr != nil {
		return Check, newRequestErr
	}
	request.Header.Set("Content-Type", "application/json")
	response, HttpDoError := http.DefaultClient.Do(request)
	if HttpDoError != nil {
		return Check, HttpDoError
	}
	resp, ReadBodyError := ioutil.ReadAll(response.Body)
	if ReadBodyError != nil {
		return Check, ReadBodyError
	}
	if !gjson.Valid(string(resp)){
		return Check, errors.New("Invalid Json")
	}
	id, _ := strconv.Atoi(gjson.Get(string(resp), "id").String())
	logger.Logger().Info("create successfully")
	return NewID(id), nil
}

func (b *Backfill)SendWechat(message string) error {
	data := map[string]interface{}{}
	mentioned_list := []string{"@all"}
	text := make(map[string]interface{})
	text["content"] = message
	text["mentioned_list"] = mentioned_list
	data["msgtype"] = "text"
	data["text"] = text
	body, _ := json.Marshal(data)
	response, err1 := http.Post(cfg.WEBHOOK_URL, "text", bytes.NewReader(body))
	if err1 != nil {
		return err1
	}
	if response.StatusCode != 200 {
		return errors.New("send wechat failed: " + message)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return nil
}

