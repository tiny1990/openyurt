package kubernetes

import (
	"errors"
	"fmt"
	"sync"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"

	"github.com/alibaba/openyurt/pkg/yurtctl/constants"
	tmplutil "github.com/alibaba/openyurt/pkg/yurtctl/util/templates"
)

const (
	ConvertJobNameBase = "yurtctl-servant-convert"
	RevertJobNameBase  = "yurtctl-servant-revert"
)

var (
	PropagationPolicy     = metav1.DeletePropagationForeground
	WaitServantJobTimeout = time.Minute * 2
	CheckServantJobPeriod = time.Second * 10
)

// YamlToObject deserializes object in yaml format to a runtime.Object
func YamlToObject(yamlContent []byte) (runtime.Object, error) {
	decode := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer().Decode
	obj, _, err := decode(yamlContent, nil, nil)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// LabelNode add a new label (<key>=<val>) to the given node
func LabelNode(cliSet *kubernetes.Clientset, node *v1.Node, key, val string) (*v1.Node, error) {
	node.Labels[key] = val
	newNode, err := cliSet.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, err
	}
	return newNode, nil
}

// AnnotateNode add a new annotation (<key>=<val>) to the given node
func AnnotateNode(cliSet *kubernetes.Clientset, node *v1.Node, key, val string) (*v1.Node, error) {
	node.Annotations[key] = val
	newNode, err := cliSet.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, err
	}
	return newNode, nil
}

// RunJobAndCleanup runs the job, wait for it to be complete, and delete it
func RunJobAndCleanup(cliSet *kubernetes.Clientset, job *batchv1.Job, timeout, period time.Duration) error {
	job, err := cliSet.BatchV1().Jobs(job.GetNamespace()).Create(job)
	if err != nil {
		return err
	}
	waitJobTimeout := time.After(timeout)
	for {
		select {
		case <-waitJobTimeout:
			return errors.New("wait for job to be complete timeout")
		case <-time.After(period):
			job, err := cliSet.BatchV1().Jobs(job.GetNamespace()).
				Get(job.GetName(), metav1.GetOptions{})
			if err != nil {
				klog.Error("fail to get job(%s) when waiting for it to be succeeded: %s",
					job.GetName(), err)
				return err
			}
			if job.Status.Succeeded == *job.Spec.Completions {
				if err := cliSet.BatchV1().Jobs(job.GetNamespace()).
					Delete(job.GetName(), &metav1.DeleteOptions{
						PropagationPolicy: &PropagationPolicy,
					}); err != nil {
					klog.Errorf("fail to delete succeeded servant job(%s): %s",
						job.GetName(), err)
					return err
				}
				return nil
			}
			continue
		}
	}
}

// RunServantJobs launchs servant jobs on specified edge nodes
func RunServantJobs(cliSet *kubernetes.Clientset, tmplCtx map[string]string, edgeNodeNames []string) error {
	var wg sync.WaitGroup
	for _, nodeName := range edgeNodeNames {
		action, exist := tmplCtx["action"]
		if !exist {
			return errors.New("action is not specified")
		}

		switch action {
		case "convert":
			tmplCtx["jobName"] = ConvertJobNameBase + "-" + nodeName
		case "revert":
			tmplCtx["jobName"] = RevertJobNameBase + "-" + nodeName
		default:
			return fmt.Errorf("unknown action: %s", action)
		}
		tmplCtx["nodeName"] = nodeName

		jobYaml, err := tmplutil.SubsituteTemplate(constants.ServantJobTemplate, tmplCtx)
		if err != nil {
			return err
		}
		srvJobObj, err := YamlToObject([]byte(jobYaml))
		if err != nil {
			return err
		}
		srvJob, ok := srvJobObj.(*batchv1.Job)
		if !ok {
			return errors.New("fail to assert yurtctl-servant job")
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := RunJobAndCleanup(cliSet, srvJob,
				WaitServantJobTimeout, CheckServantJobPeriod); err != nil {
				klog.Errorf("fail to run servant job(%s): %s",
					srvJob.GetName(), err)
			} else {
				klog.Infof("servant job(%s) has succeeded", srvJob.GetName())
			}
		}()
	}
	wg.Wait()
	return nil
}
