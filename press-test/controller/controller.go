package controller

import (
	"errors"
	cfg "press-test/config"
	logger "press-test/log"
	"press-test/utils"
)

var (
	err error
)

//func mainn() {
//	//cmd := "kubectl exec -it -n dops-onsite kafka-0 -- kafka-consumer-groups --bootstrap-server localhost:9092 --group  velocity --describe | grep velocity.gaohuitest | awk '{{if ($5 != 0) {{print $5}}}}' | wc -l '''\n"
//	cmd := "ls -ltr /Users/koko/Desktop"
//	command := exec.Command("bash", "-c", cmd)
//	output, _ := command.Output()
//	fmt.Println(string(output))
//}

func Prepare() error {
	if err = logger.Logger().InitLogFile(); err != nil {
		return err
	}
	logger.Logger().Info("----------> Start Press-test")
	if err = utils.CreateFile(); err != nil {
		return err
	}
	if err = Run(); err != nil {
		return err
	}
	return nil
}

func Run() error {
	press := utils.NewPress()
	if !press.DeletePV("cassandra") {
		logger.Logger().Error("Failed to run press-test")
		return errors.New("Failed to run press-test")
	}
	fpResource := utils.Resource{
		ControllerName: cfg.FP,
		ControllerType: cfg.DEPLOY,
		Replicas: cfg.TotalFPNum,
		Namespace: cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "",
			Limit: "",
		},
		Memory: utils.Memory{
			Request: "",
			Limit: "",
		},
	}
	dbResource := utils.Resource{
		ControllerName: cfg.CASSANDRA,
	    ControllerType: cfg.STS,
	    Replicas: cfg.TotalCassandraNum,
	    Namespace: cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "",
			Limit: "",
		},
		Memory: utils.Memory{
			Request: "",
			Limit: "",
		},
	}
	if err = press.FeatureAndDataRangeWithUpdate(&fpResource, &dbResource, cfg.IsUpdate, cfg.ConnectNumDefault); err != nil {
		return err
	}

	return nil
}
