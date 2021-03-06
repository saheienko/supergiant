package kubelet

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/util"
)

const (
	StepName = "kubelet"

	// LabelNodeRole specifies the role of a node
	LabelNodeRole = "kubernetes.io/role"
)

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (t *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.KubeletConfig.PrivateIP = config.Node.PrivateIp
	config.KubeletConfig.PublicIP = config.Node.PublicIp

	config.KubeletConfig.CACert = config.CertificatesConfig.CACert
	config.KubeletConfig.CAKey = config.CertificatesConfig.CAKey

	config.KubeletConfig.AdminKey = config.CertificatesConfig.AdminKey
	config.KubeletConfig.AdminCert = config.CertificatesConfig.AdminCert
	config.KubeletConfig.LoadBalancerHost = config.InternalDNSName
	config.KubeletConfig.NodeName = config.Node.Name
	config.KubeletConfig.IsMaster = config.IsMaster
	config.KubeletConfig.UserName = config.Kube.SSHConfig.User
	config.KubeletConfig.APIServerPort = config.Kube.APIServerPort

	if len(config.KubeletConfig.ServicesCIDR) > 0 {
		kubeDefaultSvcIp, err := util.GetKubernetesDefaultSvcIP(config.KubeletConfig.ServicesCIDR)
		if err != nil {
			return errors.Wrapf(err, "get cluster dns ip from the %s subnet", config.KubeletConfig.ServicesCIDR)
		}
		config.KubeletConfig.KubernetesSvcIP = kubeDefaultSvcIp.String()
	}

	err := steps.RunTemplate(ctx, t.script, config.Runner, out, config.KubeletConfig)

	if err != nil {
		return errors.Wrap(err, "install kubelet step")
	}

	return nil
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) Name() string {
	return StepName
}

func (t *Step) Description() string {
	return "Run kubelet"
}

func (s *Step) Depends() []string {
	return []string{docker.StepName}
}
