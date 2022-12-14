package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	cfg "press-test/config"
)

var (
	controller     = []string{"Deployment, StatefulSet", "Job", "DaemonSet", "CronJob"}
	config         *rest.Config
	client         dynamic.Interface
	result         *unstructured.Unstructured
	controllerType schema.GroupVersionResource
)

type Resource struct {
	ControllerName string
	ControllerType string
	Replicas       int
	Namespace      string
	CPU
	Memory
}

type CPU struct {
	Request string
	Limit   string
}

type Memory struct {
	Request string
	Limit   string
}

type Execute interface {
	ControllerValid(controllerType *string) bool
	ScaleReplicas(kubeconfig string) error
	ChangeCPUAndMemory(kubeconfig string) error
}

func (r *Resource) ControllerValid(controllerType *string) bool {
	for _, val := range controller {
		if val == *controllerType {
			return true
		}
	}
	return false
}

func (r *Resource) ScaleReplicas(kubeconfig string) error {
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to build kubeconfig")
		err = GenericFactory().SendWechat("failed to build kubeconfig")
		if err != nil {
			return err
		}
		return err
	}
	client, err = dynamic.NewForConfig(config)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to create k8s client")
		err = GenericFactory().SendWechat("failed to create k8s client")
		if err != nil {
			return err
		}
		return err
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	statefulsetRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}
	if r.ControllerType == "deploy" {
		controllerType = deploymentRes
	}
	if r.ControllerType == "sts" {
		controllerType = statefulsetRes
	}
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = client.Resource(controllerType).Namespace(r.Namespace).Get(context.TODO(), r.ControllerName, metav1.GetOptions{})
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to get k8s resource")
			err = GenericFactory().SendWechat("failed to get k8s resource")
			if err != nil {
				return err
			}
			return err
		}
		if err = unstructured.SetNestedField(result.Object, int64(r.Replicas), "spec", "replicas"); err != nil {
			Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to set replicas")
			err = GenericFactory().SendWechat("failed to set replicas")
			if err != nil {
				return err
			}
			return err
		}
		_, err = client.Resource(controllerType).Namespace(r.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return err
	})
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to update replicas")
		err = GenericFactory().SendWechat("failed to update replicas")
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r *Resource) ChangeCPUAndMemory(kubeconfig string) error {
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to build kubeconfig")
		err = GenericFactory().SendWechat("failed to build kubeconfig")
		if err != nil {
			return err
		}
		return err
	}
	client, err = dynamic.NewForConfig(config)
	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to create k8s client")
		err = GenericFactory().SendWechat("failed to build kubeconfig")
		if err != nil {
			return err
		}
		return err
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	statefulsetRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}
	if r.ControllerType == "Deployment" {
		controllerType = deploymentRes
	}
	if r.ControllerType == "StatefulSet" {
		controllerType = statefulsetRes
	}
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = client.Resource(controllerType).Namespace(r.Namespace).Get(context.TODO(), r.ControllerName, metav1.GetOptions{})
		if err != nil {
			Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to get k8s resource")
			err = GenericFactory().SendWechat("failed to get k8s resource")
			if err != nil {
				return err
			}
			return err
		}
		containers, found, err1 := unstructured.NestedSlice(result.Object, "spec", "template", "spec", "containers")
		if err1 != nil || !found || containers == nil {
			Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to get nestedSlice")
			err = GenericFactory().SendWechat("failed to get nestedSlice")
			if err != nil {
				return err
			}
			return err1
		}
		if r.Memory.Request != "" {
			if err = unstructured.SetNestedField(containers[0].(map[string]interface{}), r.Memory.Request, "resources", "memory", "requests"); err != nil {
				Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to set memory-request nestedSlice")
				err = GenericFactory().SendWechat("failed to set memory-request nestedSlice")
				if err != nil {
					return err
				}
				return err
			}
		}
		if r.Memory.Limit != "" {
			if err = unstructured.SetNestedField(containers[0].(map[string]interface{}), r.Memory.Limit, "resources", "memory", "limits"); err != nil {
				Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to set memory-limit nestedSlice")
				err = GenericFactory().SendWechat("failed to set memory-limit nestedSlice")
				if err != nil {
					return err
				}
				return err
			}
		}
		if r.CPU.Request != "" {
			if err = unstructured.SetNestedField(containers[0].(map[string]interface{}), r.CPU.Request, "resources", "cpu", "requests"); err != nil {
				Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to set cpu-request nestedSlice")
				err = GenericFactory().SendWechat("failed to set cpu-request nestedSlice")
				if err != nil {
					return err
				}
				return err
			}
		}
		if r.CPU.Limit != "" {
			if err = unstructured.SetNestedField(containers[0].(map[string]interface{}), r.Memory.Limit, "resources", "cpu", "limits"); err != nil {
				Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to set cpu-limit nestedSlice")
				err = GenericFactory().SendWechat("failed to set cpu-limit nestedSlice")
				if err != nil {
					return err
				}
				return err
			}
		}
		_, err = client.Resource(controllerType).Namespace(r.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return err
	})

	if err != nil {
		Log.WithFields(logrus.Fields{"tenant": cfg.Client}).Error("failed to update resources")
		err = GenericFactory().SendWechat("failed to update resources")
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
