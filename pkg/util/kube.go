package util

import (
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/workflows/steps"
)

func UpdateKubeWithCloudSpecificData(k *model.Kube, config *steps.Config) {
	logrus.Debugf("Update cloud specific data for kube %s",
		config.ClusterID)

	cloudSpecificSettings := make(map[string]string)
	logrus.Infof("Save internal DNS name %s and external DNS name %s",
		config.InternalDNSName, config.ExternalDNSName)
	k.ExternalDNSName = config.ExternalDNSName
	k.InternalDNSName = config.InternalDNSName
	k.BootstrapToken = config.BootstrapToken
	k.UserData = config.ConfigMap.Data
	k.K8SVersion = config.K8SVersion
	k.Auth.CACertHash = config.CertificatesConfig.CACertHash
	k.Auth.CertificateKey = config.KubeadmConfig.CertificateKey
	k.Auth.CACertHash = config.CertificatesConfig.CACertHash

	// Save cloudSpecificData in kube
	switch config.Provider {
	case clouds.AWS:
		// Save az to subnets mapping for this cluster
		k.Subnets = config.AWSConfig.Subnets

		// Copy data got from pre provision step to cloud specific settings of kube
		cloudSpecificSettings[clouds.AwsAZ] = config.AWSConfig.AvailabilityZone
		cloudSpecificSettings[clouds.AwsVpcCIDR] = config.AWSConfig.VPCCIDR
		cloudSpecificSettings[clouds.AwsVpcID] = config.AWSConfig.VPCID
		cloudSpecificSettings[clouds.AwsKeyPairName] = config.AWSConfig.KeyPairName
		cloudSpecificSettings[clouds.AwsMastersSecGroupID] =
			config.AWSConfig.MastersSecurityGroupID
		cloudSpecificSettings[clouds.AwsNodesSecgroupID] =
			config.AWSConfig.NodesSecurityGroupID
		// TODO(stgleb): this must be done for all types of clouds
		cloudSpecificSettings[clouds.AwsSshBootstrapPrivateKey] =
			config.Kube.SSHConfig.BootstrapPrivateKey
		cloudSpecificSettings[clouds.AwsUserProvidedSshPublicKey] =
			config.Kube.SSHConfig.PublicKey
		cloudSpecificSettings[clouds.AwsRouteTableID] =
			config.AWSConfig.RouteTableID
		cloudSpecificSettings[clouds.AwsInternetGateWayID] =
			config.AWSConfig.InternetGatewayID
		cloudSpecificSettings[clouds.AwsMasterInstanceProfile] =
			config.AWSConfig.MastersInstanceProfile
		cloudSpecificSettings[clouds.AwsNodeInstanceProfile] =
			config.AWSConfig.NodesInstanceProfile
		cloudSpecificSettings[clouds.AwsImageID] =
			config.AWSConfig.ImageID
		cloudSpecificSettings[clouds.AwsExternalLoadBalancerName] =
			config.AWSConfig.ExternalLoadBalancerName
		cloudSpecificSettings[clouds.AwsInternalLoadBalancerName] =
			config.AWSConfig.InternalLoadBalancerName
		cloudSpecificSettings[clouds.AwsVolumeSize] =
			config.AWSConfig.VolumeSize
	case clouds.GCE:
		k.Subnets = config.GCEConfig.AZs
		cloudSpecificSettings[clouds.GCETargetPoolName] = config.GCEConfig.TargetPoolName
		cloudSpecificSettings[clouds.GCEHealthCheckName] = config.GCEConfig.HealthCheckName

		cloudSpecificSettings[clouds.GCEExternalIPAddressName] = config.GCEConfig.ExternalAddressName
		cloudSpecificSettings[clouds.GCEExternalIPAddress] = config.GCEConfig.ExternalIPAddressLink
		cloudSpecificSettings[clouds.GCEExternalForwardingRuleName] = config.GCEConfig.ExternalForwardingRuleName

		cloudSpecificSettings[clouds.GCEBackendServiceLink] = config.GCEConfig.BackendServiceLink
		cloudSpecificSettings[clouds.GCEBackendServiceName] = config.GCEConfig.BackendServiceName

		cloudSpecificSettings[clouds.GCEInternalIPAddress] = config.GCEConfig.InternalIPAddressLink
		cloudSpecificSettings[clouds.GCEInternalIPAddressName] = config.GCEConfig.InternalAddressName
		cloudSpecificSettings[clouds.GCEInternalForwardingRuleName] = config.GCEConfig.InternalForwardingRuleName

		cloudSpecificSettings[clouds.GCENetworkName] = config.GCEConfig.NetworkName
		cloudSpecificSettings[clouds.GCENetworkLink] = config.GCEConfig.NetworkLink

		cloudSpecificSettings[clouds.GCEImageFamily] = config.GCEConfig.ImageFamily
	case clouds.DigitalOcean:
		cloudSpecificSettings[clouds.DigitalOceanExternalLoadBalancerID] = config.DigitalOceanConfig.ExternalLoadBalancerID
		cloudSpecificSettings[clouds.DigitalOceanInternalLoadBalancerID] = config.DigitalOceanConfig.InternalLoadBalancerID
	case clouds.Azure:
		cloudSpecificSettings[clouds.AzureVNetCIDR] = config.AzureConfig.VNetCIDR
		cloudSpecificSettings[clouds.AzureVolumeSize] = config.AzureConfig.VolumeSize
	}

	k.CloudSpec = cloudSpecificSettings
}
