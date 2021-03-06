package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
)

const CreateForwardingRulesStepName = "gce_create_forwarding_rules"

type CreateForwardingRules struct {
	timeout      time.Duration
	attemptCount int

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateForwardingRulesStep() *CreateForwardingRules {
	return &CreateForwardingRules{
		timeout:      time.Second * 10,
		attemptCount: 10,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				insertForwardingRule: func(ctx context.Context, config steps.GCEConfig, rule *compute.ForwardingRule) (*compute.Operation, error) {
					return client.ForwardingRules.Insert(config.ServiceAccount.ProjectID, config.Region, rule).Do()
				},
				getForwardingRule: func(ctx context.Context, config steps.GCEConfig, name string) (*compute.ForwardingRule, error) {
					return client.ForwardingRules.Get(config.ServiceAccount.ProjectID, config.Region, name).Do()
				},
			}, nil
		},
	}
}

func (s *CreateForwardingRules) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateForwardingRulesStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateForwardingRulesStepName)
	}

	exName := fmt.Sprintf("exrule-%s", config.ClusterID)
	externalForwardingRule := &compute.ForwardingRule{
		Name:                exName,
		IPAddress:           config.GCEConfig.ExternalIPAddressLink,
		LoadBalancingScheme: "EXTERNAL",
		Description:         "External forwarding rule to target pool",
		IPProtocol:          "TCP",
		Target:              config.GCEConfig.TargetPoolLink,
	}

	timeout := s.timeout

	for i := 0; i < s.attemptCount; i++ {
		_, err = svc.insertForwardingRule(ctx, config.GCEConfig, externalForwardingRule)

		if err == nil {
			break
		}

		logrus.Debugf("Error external forwarding rule %v sleep for %v", err, timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error creating external forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating external forwarding rule caused", CreateForwardingRulesStepName)
	}

	externalForwardingRule, err = svc.getForwardingRule(ctx, config.GCEConfig, exName)

	if err != nil {
		logrus.Errorf("get external forwarding rule %v", err)
		return errors.Wrapf(err, "get external forwarding rule")
	}

	logrus.Debugf("Created external forwarding rule %s link %s", exName, externalForwardingRule.SelfLink)
	config.GCEConfig.ExternalForwardingRuleName = externalForwardingRule.Name

	inName := fmt.Sprintf("inrule-%s", config.ClusterID)
	internalForwardingRule := &compute.ForwardingRule{
		Name:                inName,
		IPAddress:           config.GCEConfig.InternalIPAddressLink,
		LoadBalancingScheme: "INTERNAL",
		Description:         "Internal forwarding rule to target pool",
		IPProtocol:          "TCP",
		Ports:               []string{fmt.Sprintf("%d", config.Kube.APIServerPort)},
		BackendService:      config.GCEConfig.BackendServiceLink,
		Network:             config.GCEConfig.NetworkLink,
		Subnetwork:          config.GCEConfig.SubnetLink,
	}

	timeout = s.timeout

	for i := 0; i < s.attemptCount; i++ {
		_, err = svc.insertForwardingRule(ctx, config.GCEConfig, internalForwardingRule)

		if err == nil {
			break
		}

		logrus.Debugf("Error internal forwarding rule error %v sleep for %v", err, timeout)
		time.Sleep(timeout)
		timeout = timeout * 2
	}

	if err != nil {
		logrus.Errorf("Error creating internal forwarding rule %v", err)
		return errors.Wrapf(err, "%s creating internal forwarding rule caused", CreateForwardingRulesStepName)
	}

	internalForwardingRule, err = svc.getForwardingRule(ctx, config.GCEConfig, inName)

	if err != nil {
		logrus.Errorf("get internal forwarding rule %v", err)
		return errors.Wrapf(err, "get internal forwarding rule")
	}

	logrus.Debugf("Created internal forwarding rule %s link %s", inName, internalForwardingRule.SelfLink)
	config.GCEConfig.InternalForwardingRuleName = internalForwardingRule.Name

	return nil
}

func (s *CreateForwardingRules) Name() string {
	return CreateForwardingRulesStepName
}

func (s *CreateForwardingRules) Depends() []string {
	return []string{CreateTargetPullStepName, CreateIPAddressStepName}
}

func (s *CreateForwardingRules) Description() string {
	return "Create forwarding rules to pass traffic to nodes"
}

func (s *CreateForwardingRules) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
