package access

import (
	"kingfisher/kf/common/access"
	"kingfisher/kf/common/log"
	authenticationv1alpha1 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/authentication/v1alpha1"
	configv1alpha2 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/config/v1alpha2"
	networkingv1alpha3 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	securityv1beta1 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/security/v1beta1"
)

func IstioNetworkingClient(clusterId string) (*networkingv1alpha3.NetworkingV1alpha3Client, error) {
	config, err := access.GetConfig(clusterId)
	if err != nil {
		return nil, err
	}
	clientSet, err := networkingv1alpha3.NewForConfig(config)
	if err != nil {
		log.Errorf("Istio Networking clientSet error: %s", err)
		return nil, err
	}
	return clientSet, err
}

func IstioAuthenticationClient(clusterId string) (*authenticationv1alpha1.AuthenticationV1alpha1Client, error) {
	config, err := access.GetConfig(clusterId)
	if err != nil {
		return nil, err
	}
	clientSet, err := authenticationv1alpha1.NewForConfig(config)
	if err != nil {
		log.Errorf("Istio Authentication clientSet error: %s", err)
		return nil, err
	}
	return clientSet, err
}

func IstioSecurityClient(clusterId string) (*securityv1beta1.SecurityV1beta1Client, error) {
	config, err := access.GetConfig(clusterId)
	if err != nil {
		return nil, err
	}
	clientSet, err := securityv1beta1.NewForConfig(config)
	if err != nil {
		log.Errorf("Istio Security clientSet error: %s", err)
		return nil, err
	}
	return clientSet, err
}

func IstioConfigClient(clusterId string) (*configv1alpha2.ConfigV1alpha2Client, error) {
	config, err := access.GetConfig(clusterId)
	if err != nil {
		return nil, err
	}
	clientSet, err := configv1alpha2.NewForConfig(config)
	if err != nil {
		log.Errorf("Istio Config clientSet error: %s", err)
		return nil, err
	}
	return clientSet, err
}
