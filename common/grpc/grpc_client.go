package grpc

import (
	"context"
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"kingfisher/kf/common"
	"kingfisher/kf/common/grpc"
	"kingfisher/kf/common/log"
	pb "kingfisher/king-k8s/grpc/proto"
	"time"
)

func GetService(cluster, namespace, name string) (*v1.Service, error) {
	conn := grpc.ClientDial(common.KingK8S)
	defer conn.Close()
	c := pb.NewServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	response, err := c.Get(ctx, &pb.ServiceRequest{Cluster: cluster, Namespace: namespace, Name: name})
	if err != nil {
		log.Errorf("could not grpc request: %v", err)
		return nil, err
	}
	service := &v1.Service{}
	err = json.Unmarshal(response.Data, service)
	if err != nil {
		log.Errorf("grpc response unmarshal error: %v", err)
		return nil, err
	}
	return service, nil
}

func GetDeployment(cluster, namespace, labels string) (*appsv1.DeploymentList, error) {
	// get deployment
	conn := grpc.ClientDial(common.KingK8S)
	defer conn.Close()
	c := pb.NewDeploymentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	response, err := c.GetByLabels(ctx, &pb.DeploymentRequest{Cluster: cluster, Labels: labels, Namespace: namespace})
	if err != nil {
		log.Errorf("could not grpc request: %v", err)
		return nil, err
	}
	deployment := &appsv1.DeploymentList{}
	err = json.Unmarshal(response.Data, deployment)
	if err != nil {
		log.Errorf("grpc response unmarshal error: %v", err)
		return nil, err
	}
	return deployment, nil
}
