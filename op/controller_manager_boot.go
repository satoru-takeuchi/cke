package op

import (
	"context"

	"github.com/cybozu-go/cke"
	"github.com/cybozu-go/cke/common"
	"k8s.io/client-go/tools/clientcmd"
)

type controllerManagerBootOp struct {
	nodes []*cke.Node

	cluster       string
	serviceSubnet string
	params        cke.ServiceParams

	step  int
	files *common.FilesBuilder
}

// ControllerManagerBootOp returns an Operator to bootstrap kube-controller-manager
func ControllerManagerBootOp(nodes []*cke.Node, cluster string, serviceSubnet string, params cke.ServiceParams) cke.Operator {
	return &controllerManagerBootOp{
		nodes:         nodes,
		cluster:       cluster,
		serviceSubnet: serviceSubnet,
		params:        params,
		files:         common.NewFilesBuilder(nodes),
	}
}

func (o *controllerManagerBootOp) Name() string {
	return "kube-controller-manager-bootstrap"
}

func (o *controllerManagerBootOp) NextCommand() cke.Commander {
	switch o.step {
	case 0:
		o.step++
		return common.ImagePullCommand(o.nodes, cke.HyperkubeImage)
	case 1:
		o.step++
		dirs := []string{
			"/var/log/kubernetes/controller-manager",
		}
		return common.MakeDirsCommand(o.nodes, dirs)
	case 2:
		o.step++
		return prepareControllerManagerFilesCommand{o.cluster, o.files}
	case 3:
		o.step++
		return o.files
	case 4:
		o.step++
		return common.RunContainerCommand(o.nodes,
			kubeControllerManagerContainerName, cke.HyperkubeImage,
			common.WithParams(ControllerManagerParams(o.cluster, o.serviceSubnet)),
			common.WithExtra(o.params))
	default:
		return nil
	}
}

type prepareControllerManagerFilesCommand struct {
	cluster string
	files   *common.FilesBuilder
}

func (c prepareControllerManagerFilesCommand) Run(ctx context.Context, inf cke.Infrastructure) error {
	const kubeconfigPath = "/etc/kubernetes/controller-manager/kubeconfig"
	storage := inf.Storage()

	ca, err := storage.GetCACertificate(ctx, "kubernetes")
	if err != nil {
		return err
	}
	g := func(ctx context.Context, n *cke.Node) ([]byte, error) {
		crt, key, err := cke.KubernetesCA{}.IssueForControllerManager(ctx, inf)
		if err != nil {
			return nil, err
		}
		cfg := controllerManagerKubeconfig(c.cluster, ca, crt, key)
		return clientcmd.Write(*cfg)
	}
	err = c.files.AddFile(ctx, kubeconfigPath, g)
	if err != nil {
		return err
	}

	saKey, err := storage.GetServiceAccountKey(ctx)
	if err != nil {
		return err
	}
	saKeyData := []byte(saKey)
	g = func(ctx context.Context, n *cke.Node) ([]byte, error) {
		return saKeyData, nil
	}
	return c.files.AddFile(ctx, K8sPKIPath("service-account.key"), g)
}

func (c prepareControllerManagerFilesCommand) Command() cke.Command {
	return cke.Command{
		Name: "prepare-controller-manager-files",
	}
}

// ControllerManagerParams returns parameters for kube-controller-manager.
func ControllerManagerParams(clusterName, serviceSubnet string) cke.ServiceParams {
	args := []string{
		"controller-manager",
		"--cluster-name=" + clusterName,
		"--service-cluster-ip-range=" + serviceSubnet,
		"--kubeconfig=/etc/kubernetes/controller-manager/kubeconfig",
		"--log-dir=/var/log/kubernetes/controller-manager",
		"--logtostderr=false",

		// ToDo: cluster signing
		// https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#a-note-to-cluster-administrators
		// https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet-tls-bootstrapping/
		//    Create an intermediate CA under cke/ca-kubernetes?
		//    or just an certficate/key pair?
		// "--cluster-signing-cert-file=..."
		// "--cluster-signing-key-file=..."

		// for service accounts
		"--root-ca-file=" + K8sPKIPath("ca.crt"),
		"--service-account-private-key-file=" + K8sPKIPath("service-account.key"),
		"--use-service-account-credentials=true",
	}
	return cke.ServiceParams{
		ExtraArguments: args,
		ExtraBinds: []cke.Mount{
			{"/etc/machine-id", "/etc/machine-id", true, "", ""},
			{"/etc/kubernetes", "/etc/kubernetes", true, "", cke.LabelShared},
			{"/var/log/kubernetes/controller-manager", "/var/log/kubernetes/controller-manager", false, "", cke.LabelPrivate},
		},
	}
}