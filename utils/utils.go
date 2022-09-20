package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RealGaohui/urlBuilder"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/exec"
	"path"
	cfg "press-test/config"
	"strconv"
	"strings"
	"time"
)

var (
	command *exec.Cmd
	err     error
	output  []byte
)

type Press struct {
	TenantContext context.Context
	SSHClient     *ssh.Client
}

type fpTenant struct {
	Name string `json:"name"`
}

type execute interface {
	CreateOrDeletePvPath(tenantContext context.Context, path, host, action string) error
	CreateFPTanant(tenantContext context.Context) (bool, error)
	WaitKafkaConsumer(tenantContext context.Context) bool
	DeletePV(tenantContext context.Context, db string) bool
	PrepareDir(tenantContext context.Context, path string) error
	CreateFile(tenantContext context.Context, dir, name string) error
	RecordWrkLog(tenantContext context.Context, msg string) error
	GenerateWrkResult(tenantContext context.Context, wrkResult string) (*WRK, error)
	FeatureAndDataRangeWithUpdate(tenantContext context.Context, fpResource, dbResource *Resource, isUpdate bool, connectNum int) error
	FeatureAndDataRangeWithoutUpdate(tenantContext context.Context, fpResource, dbResource *Resource, isUpdate bool, connectNum int) error
	EndpointIncreasing(tenantContext context.Context, resource *Resource)
	FPCPUIncreasing(tenantContext context.Context, cpu *CPU)
	CassandraCPUIncreasing(tenantContext context.Context, cpu *CPU)
	CopyTemplate(tenantContext context.Context, template, index string) (string, error)
	RestartFP(tenantContext context.Context, namespace string) bool
	WaitForReady(tenantContext context.Context, deployName, deployType, namespace string) bool
	WriteCSVHeader(tenantContext context.Context) error
}

func NewPress(tenantContext context.Context) execute {
	return &Press{
		TenantContext: tenantContext,
		SSHClient:     sshClient,
	}
}

func (p *Press) CreateOrDeletePvPath(tenantContext context.Context, path, host, action string) error {
	var ip IP
	ip = IP(host)
	sshclient, sshInitError := ip.SSHConnect()
	if sshInitError != nil {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to create ssh client")
		return sshInitError
	}
	session, newSessionError := sshclient.NewSession()
	if newSessionError != nil {
		return newSessionError
	}
	if action == cfg.ACTIONDELETE {
		_, err = session.CombinedOutput(fmt.Sprintf("rm -rf %s", path))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error(fmt.Sprintf("failed to %s %s:%s", action, host, path))
			return err
		}
	}
	if action == cfg.ACTIONCREATE {
		_, err = session.CombinedOutput(fmt.Sprintf("mkdir -p %s", path))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error(fmt.Sprintf("failed to %s %s:%s", action, host, path))
			return err
		}
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info(fmt.Sprintf("%s %s:%s successfully", action, host, path))
	return nil
}

func (p *Press) CreateFPTanant(tenantContext context.Context) (bool, error) {
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("create fp tenant for: [%s]", tenantContext.Value("tenant"))
	URL := urlBuilder.URLBuilder().
		SetBase(cfg.FpBaseUrl).
		SetPath("/" + tenantContext.Value("tenant").(string)).
		SetPath("/config/tenant/create").
		ToString()
	FPTenant := fpTenant{
		tenantContext.Value("tenant").(string),
	}
	body, _ := json.Marshal(FPTenant)
	request, newRequestErr := http.NewRequest("POST", URL, bytes.NewReader(body))
	if newRequestErr != nil {
		return false, newRequestErr
	}
	request.Header.Set("Content-Type", "application/json")
	response, HttpDoError := http.DefaultClient.Do(request)
	if HttpDoError != nil {
		return false, HttpDoError
	}
	if response.StatusCode != 200 {
		return false, errors.New("http error")
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("create fp tenant successfully for [%s]", tenantContext.Value("tenant"))
	return true, nil
}

func (p *Press) WaitKafkaConsumer(tenantContext context.Context) bool {
	for {
		cmd := fmt.Sprintf("kubectl --kubeconfig=%s exec -it -n %s kafka-0 -- kafka-consumer-groups --bootstrap-server localhost:9092 --group  velocity --describe | grep velocity.%s | awk '{{if ($5 != 0) {{print $5}}}}' | wc -l '''\n", cfg.KubeconfigPATH, cfg.K8sNamespaceCassandra, tenantContext.Value("tenant"))
		waitKafkaConsumerResult, waitKafkaConsumerCommand := GenericFactory().Command(cmd)
		if waitKafkaConsumerCommand != nil {
			return false
		}
		tmp := strconv.Itoa(0)
		if waitKafkaConsumerResult == tmp {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("kafka consumer is  0, press-test completed successfully")
			return true
		}
	}
}

func (p *Press) PrepareDir(tenantContext context.Context, path string) error {
	if !checkExist(tenantContext, path) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Press) CreateFile(tenantContext context.Context, dir, name string) error {
	err = p.PrepareDir(tenantContext, dir)
	if err != nil {
		return err
	}
	file := path.Join(dir, name)
	if !checkExist(tenantContext, file) {
		_, err = os.Create(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkExist(tenantContext context.Context, path string) bool {
	_, err = os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (p *Press) DeletePV(tenantContext context.Context, db string) bool {
	//get themselves prefix namespace grepkey for each db
	prefix := ""
	namespace := ""
	grepKey := ""
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("start delete pv %s", db)
	if db == cfg.CASSANDRA {
		prefix = cfg.CASSANDRAPREFIX
		namespace = cfg.K8sNamespaceCassandra
		grepKey = db
	}
	if db == cfg.YB {
		prefix = cfg.YBPREFIX
		namespace = cfg.K8sNamespace
		grepKey = "yb-"
	}
	//get po by kubectl
	//pod slice
	pods := ""
	Pods := make([]string, 3)
	if db == cfg.CASSANDRA {
		pods, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s get po | grep %s | grep -v cassandra-tool | awk '{print $1}'", cfg.KubeconfigPATH, namespace, grepKey))
		podSlice := strings.Split(pods, "\n")
		Pods = podSlice[:len(podSlice)-1]
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to get pods")
			return false
		}
	} else {
		pods, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s get po | grep %s | awk '{print $1}'", cfg.KubeconfigPATH, namespace, grepKey))
		podSlice := strings.Split(pods, "\n")
		Pods = podSlice[:len(podSlice)-1]
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to get pods")
			return false
		}
	}
	//for delete pvc and pv, scale down
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("trying to delete pv for pods %s", Pods)
	if len(Pods) == 0 {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warn("no cassandra pods, continue next")
		return true
	}
	if db == cfg.CASSANDRA {
		_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s scale statefulset %s --replicas 0", cfg.KubeconfigPATH, namespace, db))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to scale cassandra to 0")
			return false
		}
	}
	if db == cfg.YB {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("delete yb")
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("wait for cassandra to scale 0 or delete yb")
	n := 0
	for {
		getPodsResult := ""
		num := make([]string, 3)
		if db == cfg.CASSANDRA {
			getPodsResult, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s get po | grep %s | grep -v cassandra-tool | awk '{print $1}'", cfg.KubeconfigPATH, namespace, grepKey))
			podResult := strings.Split(getPodsResult, "\n")
			num = podResult[:len(podResult)-1]
			if err != nil {
				Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to get pods for waiting")
				return false
			}
		} else {
			getPodsResult, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s get po | grep %s | awk '{print $1}'", cfg.KubeconfigPATH, namespace, grepKey))
			podResult := strings.Split(getPodsResult, "\n")
			num = podResult[:len(podResult)-1]
			if err != nil {
				Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to get pods for waiting")
				return false
			}
		}
		if len(num) == 0 {
			//break loop if no cassandra pod
			break
		}
		time.Sleep(5 * time.Second)
		n++
		if n == 60 {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warn("failed to scale cassandra or delete yb in 5 mins")
		}
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("scale cassandra to 0 or delete yb completed")

	//delete pvc for each pod
	for _, pod := range Pods {
		_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s delete pvc %s", cfg.KubeconfigPATH, namespace, prefix+"-"+pod))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Errorf("failed to delete pvc [%s] for [%s]", prefix+"-"+pod, pod)
			return false
		}
		index := ""
		pvName := ""
		if strings.Split(pod, "-")[1] == "0" {
			index = "1"
			pvName = strings.Split(pod, "-")[0] + "-" + cfg.K8sNamespaceCassandra + "-" + index
		} else if strings.Split(pod, "-")[1] == "1" {
			index = "0"
			pvName = strings.Split(pod, "-")[0] + "-" + cfg.K8sNamespaceCassandra + "-" + index
		} else {
			pvName = strings.Split(pod, "-")[0] + "-" + cfg.K8sNamespaceCassandra + "-" + strings.Split(pod, "-")[1]
		}

		//delete pv
		_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s delete pv %s", cfg.KubeconfigPATH, namespace, pvName))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Errorf("failed to delete pv [%s] for [%s]", pvName, pod)
			return false
		}
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("delete pv completely for [%s]", pod)
		//delete hostpath
		hosts := strings.Split(strings.Replace(cfg.Host, " ", "", -1), ",")
		for _, host := range hosts {
			if err = p.CreateOrDeletePvPath(tenantContext, cfg.CassandraDataPath, host, cfg.ACTIONDELETE); err != nil {
				return false
			}
		}
	}
	//deploy yb
	if db == cfg.YB {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("start to deploy yb")
		_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s apply -f  yb.yaml", cfg.KubeconfigPATH, namespace))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to deploy yb")
			return false
		}
		//wait yb pod is running
		if !p.WaitForReady(tenantContext, db, cfg.STS, namespace) {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warnf("[%s] is not running", db)
			return false
		}
		//restart fp
		if !p.RestartFP(tenantContext, namespace) {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to restart fp")
			return false
		}
	}
	return true
}

func (p *Press) RestartFP(tenantContext context.Context, namespace string) bool {
	_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s rollout restart deploy fp-deployment", cfg.KubeconfigPATH, namespace))
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("failed to restart fp")
		return false
	}
	//wait for fp ready
	if !p.WaitForReady(tenantContext, cfg.FP, cfg.DEPLOY, namespace) {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warn("fp is not running")
		return false
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("fp is ready")
	return true
}

func (p *Press) WaitForReady(tenantContext context.Context, deployName, deployType, namespace string) bool {
	n := 0
	for {
		getDeployResult, getDeployErr := GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig %s -n %s get %s %s | awk '{print $2}'", cfg.KubeconfigPATH, namespace, deployType, deployName))
		if getDeployErr != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Errorf("failed to get [%s] for running", deployName)
			return false
		}
		fp := strings.Split(getDeployResult, "/")
		if fp[0] == fp[1] {
			//break loop if no cassandra pod
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("[%s] ready", deployName)
			break
		}
		time.Sleep(5 * time.Second)
		n++
		if n > 360 {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warnf("[%s] is not running in 60 minutes", deployName)
		}
	}
	//Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("%s is running", deployName)
	return true
}

func (p *Press) RecordWrkLog(tenantContext context.Context, msg string) error {
	currentTime := time.Now().Format("20060102_150405")
	fileName := cfg.WrkRawlogPath + "/" + currentTime + ".log"
	if err = initFile(fileName); err != nil {
		return err
	}
	file, openFileError := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileError != nil {
		return openFileError
	}
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(msg)
	if err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

// TODO
func (p *Press) GenerateWrkResult(tenantContext context.Context, wrkResult string) (*WRK, error) {
	generateResult := strings.Split(strings.Replace(wrkResult, "\n", " ", -1), " ")
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info(len(generateResult))
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info(generateResult)
	generateResult8, _ := strconv.Atoi(generateResult[8])
	Wrk := WRK{
		Total:       generateResult[39],
		Avg:         generateResult[18],
		Stdev:       generateResult[19],
		Max:         generateResult[20],
		P50:         generateResult[30],
		P75:         generateResult[32],
		P90:         generateResult[34],
		P95:         generateResult[36],
		P99:         generateResult[38],
		Connections: generateResult8,
	}
	if len(generateResult) == 59 {
		Wrk.QPS = generateResult[56]
		Wrk.Timeout = generateResult[54]
	} else if len(generateResult) == 54 {
		Wrk.QPS = generateResult[51]
	} else {
		Wrk.QPS = generateResult[46]
	}
	return &Wrk, nil
}

func (p *Press) CopyTemplate(tenantContext context.Context, template, index string) (string, error) {
	templateFile := strings.Split(template, ".")[0] + "-" + index + ".yaml"
	_, err = GenericFactory().Command(fmt.Sprintf("cp %s %s", template, templateFile))
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("failed to copy template file [%s]", templateFile)
		return "", err
	}
	return templateFile, nil
}

func (p *Press) FeatureAndDataRangeWithUpdate(tenantContext context.Context, fpResource, dbResource *Resource, isUpdate bool, connectNum int) error {
	//prepare db hostpath
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("prepare pv host path")
	for i := 0; i < dbResource.Replicas; i++ {
		hosts := strings.Split(strings.Replace(cfg.Host, " ", "", -1), ",")
		for _, host := range hosts {
			Path := cfg.CassandraData1Path + strconv.Itoa(i) + "-0"
			if err = p.CreateOrDeletePvPath(tenantContext, Path, host, cfg.ACTIONCREATE); err != nil {
				return err
			}
		}
		//copy pv template
		newfile, copyTempalteError := p.CopyTemplate(tenantContext, cfg.TemplatePath, strconv.Itoa(i))
		if copyTempalteError != nil {
			return copyTempalteError
		}

		//change pv template
		pvName := cfg.CASSANDRA + "-" + cfg.K8sNamespaceCassandra + "-" + strconv.Itoa(i)
		host := hosts[i]
		hostPath := cfg.CassandraDataPath + strconv.Itoa(i) + "-0"
		pv(&newfile).
			ChangePVPath(pvName, hostPath).
			ChangePVName(pvName).
			ChangePVHost(pvName, host).
			ChangeStorageClassName(pvName, cfg.K8sNamespaceCassandra)

		//create pv
		_, err = GenericFactory().Command(fmt.Sprintf("kubectl --kubeconfig=%s -n %s apply -f %s", cfg.KubeconfigPATH, cfg.K8sNamespaceCassandra, newfile))
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Warnf("failed to create pv [%s]", pvName)
			return err
		}
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof(" create pv [%s] successfully", pvName)
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("change db resource")
	//change db resource
	if err = dbResource.ChangeCPUAndMemory(cfg.KubeconfigPATH); err != nil {
		return err
	}
	if !p.WaitForReady(tenantContext, dbResource.ControllerName, dbResource.ControllerType, dbResource.Namespace) {
		return errors.New(fmt.Sprintf("[%s] is not running", dbResource.ControllerName))
	}
	//change fp resource
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("change fp resource")
	if err = fpResource.ChangeCPUAndMemory(cfg.KubeconfigPATH); err != nil {
		return err
	}
	if !p.WaitForReady(tenantContext, fpResource.ControllerName, fpResource.ControllerType, fpResource.Namespace) {
		return errors.New(fmt.Sprintf("[%s] is not running", fpResource.ControllerName))
	}
	//write csv header
	if err = p.WriteCSVHeader(tenantContext); err != nil {
		return err
	}
	//create fp tenant
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("create fp tenant")
	ok, createTenantError := p.CreateFPTanant(tenantContext)
	if !ok && createTenantError != nil {
		Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Error("Failed to create fp tenant")
		return createTenantError
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("change fp replicas")
	if err = fpResource.ScaleReplicas(cfg.KubeconfigPATH); err != nil {
		return err
	}
	if !p.WaitForReady(tenantContext, fpResource.ControllerName, fpResource.ControllerType, fpResource.Namespace) {
		return errors.New(fmt.Sprintf("[%s] is not running", fpResource.ControllerName))
	}
	//run wrk
	featureRanges := strings.Split(strings.Replace(cfg.Feature, " ", "", -1), ",")
	dataRanges := strings.Split(strings.Replace(cfg.DataRange, " ", "", -1), ",")
	for _, featureRange := range featureRanges {
		for _, dataRange := range dataRanges {
			backfillStartTime := time.Now()
			backfill := Backfill{}
			backfillOK, backfillError := backfill.DoBackfill()
			if backfillError != nil || !backfillOK.CheckBackfillFinish() {
				return backfillError
			}
			useTime := time.Since(backfillStartTime)
			title := fmt.Sprintf("%s:%s = %s:%s, connect: %s",
				fpResource.ControllerName,
				dbResource.ControllerName,
				fpResource.Replicas,
				dbResource.Replicas,
				connectNum,
			)
			//result report to wechat
			reportResult := Result{
				Backfill: useTime,
				FP: Resource{
					Replicas: fpResource.Replicas,
					CPU:      fpResource.CPU,
					Memory:   fpResource.Memory,
				},
				DB: Resource{
					Replicas: dbResource.Replicas,
					CPU:      dbResource.CPU,
					Memory:   dbResource.Memory,
				},
			}
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("data warm up")
			_, commandError := GenericFactory().Command(fmt.Sprintf("sh %s %d 2 %s %s %v", cfg.WrkScript, connectNum, featureRange, dataRange, isUpdate))
			if commandError != nil {
				return commandError
			}
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("data warm up completed")
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Infof("---------- %v ----------", title)
			Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("wait for kafka consumer to be 0")
			if !p.WaitForReady(tenantContext, cfg.KAFKA, cfg.STS, cfg.K8sNamespaceCassandra) {
				return errors.New("kafka consumer is not 0")
			}

			//start press-test
			pressTestOutput, wrkCommandError := GenericFactory().Command(fmt.Sprintf("sh %s %d 2 %s %s %v", cfg.WrkScript, connectNum, featureRange, dataRange, isUpdate))
			if wrkCommandError != nil {
				return wrkCommandError
			}
			if err = p.RecordWrkLog(tenantContext, fmt.Sprintf("%s \n %s", title, pressTestOutput)); err != nil {
				return err
			}
			wrkResult, generateWrkResultError := p.GenerateWrkResult(tenantContext, pressTestOutput)
			if generateWrkResultError != nil {
				return generateWrkResultError
			}
			//write result to csv
			if err = wrkResult.WriteToCSV(title, &reportResult); err != nil {
				return err
			}
			//report to wechat
			if err = wrkResult.Report(&reportResult); err != nil {
				return err
			}

		}
	}
	Log.WithFields(logrus.Fields{"tenant": tenantContext.Value("tenant")}).Info("press-test completed")
	return nil
}

func (p *Press) FeatureAndDataRangeWithoutUpdate(tenantContext context.Context, fpResource, dbResource *Resource, isUpdate bool, connectNum int) error {
	if err = p.FeatureAndDataRangeWithUpdate(tenantContext, fpResource, dbResource, isUpdate, connectNum); err != nil {
		return err
	}
	return nil
}

func (p *Press) WriteCSVHeader(tenantContext context.Context) error {
	file, openFileError := os.Open(cfg.CsvFilePath + "/" + cfg.CsvFileName)
	if openFileError != nil {
		return openFileError
	}
	writer := bufio.NewWriter(file)
	tmpHeader := make([]string, 10)
	tmpHeader = append(tmpHeader, "资源配置")
	tmpHeader = append(tmpHeader, "连接数")
	tmpHeader = append(tmpHeader, "qps")
	tmpHeader = append(tmpHeader, "P50")
	tmpHeader = append(tmpHeader, "P75")
	tmpHeader = append(tmpHeader, "P90")
	tmpHeader = append(tmpHeader, "P95")
	tmpHeader = append(tmpHeader, "P99")
	tmpHeader = append(tmpHeader, "avg")
	tmpHeader = append(tmpHeader, "stdev")
	tmpHeader = append(tmpHeader, "max")
	tmpHeader = append(tmpHeader, "backfill_use_time")
	_, err = writer.WriteString(strings.Join(tmpHeader, ",") + "\n")
	if err != nil {
		return err
	}
	_ = writer.Flush()
	_ = file.Close()
	return nil
}

func (p *Press) EndpointIncreasing(tenantContext context.Context, resource *Resource) {

}

func (p *Press) FPCPUIncreasing(tenantContext context.Context, cpu *CPU) {

}

func (p *Press) CassandraCPUIncreasing(tenantContext context.Context, cpu *CPU) {

}
