package controller

import (
	"golang.org/x/net/context"
	cfg "press-test/config"
	"press-test/utils"
)

var (
	err error
)

func Prepare(tenantContext context.Context) error {
	if err = utils.InitLog(tenantContext); err != nil {
		return err
	}
	if err = utils.CreateFile(); err != nil {
		return err
	}
	if err = Run(tenantContext); err != nil {
		return err
	}
	return nil
}

func Run(tenantContext context.Context) error {
	press := utils.NewPress(tenantContext)
	//if !press.DeletePV(tenantContext, cfg.CASSANDRA) {
	//	return errors.New("failed to run press-test")
	//}
	fpResource := utils.Resource{
		ControllerName: cfg.FP,
		ControllerType: cfg.DEPLOY,
		Replicas:       cfg.TotalFPNum,
		Namespace:      cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "2",
			Limit:   "2",
		},
		Memory: utils.Memory{
			Request: "4",
			Limit:   "4",
		},
	}
	dbResource := utils.Resource{
		ControllerName: cfg.CASSANDRA,
		ControllerType: cfg.STS,
		Replicas:       cfg.TotalCassandraNum,
		Namespace:      cfg.K8sNamespaceCassandra,
		CPU: utils.CPU{
			Request: "2",
			Limit:   "2",
		},
		Memory: utils.Memory{
			Request: "4",
			Limit:   "4",
		},
	}
	if err = press.FeatureAndDataRangeWithUpdate(tenantContext, &fpResource, &dbResource, cfg.IsUpdate, cfg.ConnectNumDefault); err != nil {
		return err
	}

	return nil
}
