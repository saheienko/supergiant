package cloudcontroller

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "cloudcontroller"

type Config struct {
	K8SVersion      string
	KubeadmVersion  string
	IsBootstrap     bool
	IsMaster        bool
	InternalDNSName string
	ExternalDNSName string
	Token           string
	CACertHash      string
	CertificateKey  string
	CIDR            string
	ServiceCIDR     string
	UserName        string
	Provider        string
	APIServerPort   int64
	NodeIp          string
}

type Step struct {
	script *template.Template
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	err := steps.RunTemplate(context.Background(), s.script, config.Runner, out, toStepCfg(config))

	if err != nil {
		return errors.Wrap(err, "install cloud-controller-manager")
	}

	return nil
}

func (s *Step) Rollback(ctx context.Context, out io.Writer, config *steps.Config) error {
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "create cloud-controller-manager"
}

func (s *Step) Depends() []string {
	return nil
}

func toStepCfg(c *steps.Config) Config {
	return Config{
		KubeadmVersion:  "1.15.0", // TODO(stgleb): get it from available versions once we have them
		K8SVersion:      c.Kube.K8SVersion,
		IsBootstrap:     c.IsBootstrap,
		IsMaster:        c.IsMaster,
		InternalDNSName: c.Kube.InternalDNSName,
		ExternalDNSName: c.Kube.ExternalDNSName,
		Token:           c.Kube.BootstrapToken,
		CACertHash:      c.Kube.Auth.CACertHash,
		CertificateKey:  c.Kube.Auth.CertificateKey,
		CIDR:            c.Kube.Networking.CIDR,
		ServiceCIDR:     c.Kube.ServicesCIDR,
		UserName:        clouds.OSUser,
		//Provider:        toCloudProviderOpt(c.Kube.Provider),
		APIServerPort: c.Kube.APIServerPort,
		NodeIp:        c.Node.PrivateIp,
	}
}
