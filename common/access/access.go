package access

import (
	"github.com/open-kingfisher/king-utils/common/access"
	"github.com/open-kingfisher/king-utils/common/log"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
)

func IstioClient(clusterId string) (*versionedclient.Clientset, error) {
	config, err := access.GetConfig(clusterId)
	if err != nil {
		return nil, err
	}
	clientSet, err := versionedclient.NewForConfig(config)
	if err != nil {
		log.Errorf("Istio clientSet error: %s", err)
		return nil, err
	}
	return clientSet, err
}
