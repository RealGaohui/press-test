package controller

import (
	"errors"
	cfg "press-test/config"
	"press-test/utils"
)

var (
	err error
)

func Prepare() error {
	if err = utils.InitLog(); err != nil {
		return err
	}
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
		return errors.New("failed to run press-test")
	}
	fpResource := utils.Resource{
		ControllerName: cfg.FP,
		ControllerType: cfg.DEPLOY,
		Replicas:       cfg.TotalFPNum,
		Namespace:      cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "",
			Limit:   "",
		},
		Memory: utils.Memory{
			Request: "",
			Limit:   "",
		},
	}
	dbResource := utils.Resource{
		ControllerName: cfg.CASSANDRA,
		ControllerType: cfg.STS,
		Replicas:       cfg.TotalCassandraNum,
		Namespace:      cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "",
			Limit:   "",
		},
		Memory: utils.Memory{
			Request: "",
			Limit:   "",
		},
	}
	if err = press.FeatureAndDataRangeWithUpdate(&fpResource, &dbResource, cfg.IsUpdate, cfg.ConnectNumDefault); err != nil {
		return err
	}

	return nil
}
