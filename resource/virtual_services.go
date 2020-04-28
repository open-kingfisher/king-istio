package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-istio/common/chart"
	grpcClient "github.com/open-kingfisher/king-istio/common/grpc"
	"github.com/open-kingfisher/king-k8s/util"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/kit"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type VirtualServicesResource struct {
	Params   *handle.Resources
	PostData *v1alpha3.VirtualService
	Access   *versionedclient.Clientset
}

func (r *VirtualServicesResource) Get() (*v1alpha3.VirtualService, error) {
	var ctx context.Context
	i, err := r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Get(ctx, r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "VirtualService", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *VirtualServicesResource) List() (*v1alpha3.VirtualServiceList, error) {
	var ctx context.Context
	return r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).List(ctx, metav1.ListOptions{})
}

func (r *VirtualServicesResource) Delete() (err error) {
	var ctx context.Context
	if err = r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Delete(ctx, r.Params.Name, metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.VirtualServices,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *VirtualServicesResource) Patch() (res *v1alpha3.VirtualService, err error) {
	var ctx context.Context
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Patch(ctx, r.Params.Name, types.JSONPatchType, data, metav1.PatchOptions{}); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.VirtualServices, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.VirtualServices,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *VirtualServicesResource) Update() (res *v1alpha3.VirtualService, err error) {
	var ctx context.Context
	if r.Params.PostType == "form" {
		if vs, err := r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Get(ctx, r.PostData.Name, metav1.GetOptions{}); err != nil {
			log.Errorf("%s get error:%s; Json:%+v; Name:%s", common.VirtualServices, err, r.PostData, r.PostData.Name)
			return nil, err
		} else {
			r.PostData.Spec.Tcp = vs.Spec.Tcp
			r.PostData.Spec.Tls = vs.Spec.Tls
			r.PostData.Spec.ExportTo = vs.Spec.ExportTo
		}
	}
	if res, err = r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Update(ctx, r.PostData, metav1.UpdateOptions{}); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.VirtualServices, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.VirtualServices,
		ActionType: common.Update,
		PostType:   common.ActionType(r.Params.PostType),
		Resources:  r.Params,
		Name:       r.PostData.Name,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *VirtualServicesResource) Create() (res *v1alpha3.VirtualService, err error) {
	var ctx context.Context
	if res, err = r.Access.NetworkingV1alpha3().VirtualServices(r.Params.Namespace).Create(ctx, r.PostData, metav1.CreateOptions{}); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.VirtualServices, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.VirtualServices,
		ActionType: common.Create,
		PostType:   common.ActionType(r.Params.PostType),
		Resources:  r.Params,
		Name:       r.PostData.Name,
		PostData:   r.PostData,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *VirtualServicesResource) Chart() (interface{}, error) {
	chartData := chart.Chart{}
	if vs, err := r.Get(); err != nil {
		return nil, err
	} else {
		count := 0
		chartData.Name = r.Params.Name
		chartData.Rank = "vs"
		for _, vs := range vs.Spec.Http {
			route := chart.Chart{}
			route.Name = "route"
			route.Rank = "route"
			for _, routers := range vs.Route {
				destination := chart.Chart{}
				destination.Name = fmt.Sprintf("主机: %s 权重: %d", routers.Destination.Host, routers.Weight)
				destination.Rank = "destination"
				dr := DestinationRulesResource{
					Params: r.Params,
					Access: *r.Access,
				}
				drs, _ := dr.List()
				for _, d := range drs.Items {
					if d.Spec.Host == routers.Destination.Host {
						for _, dd := range d.Spec.Subsets {
							if dd.Name == routers.Destination.Subset {
								label := make(map[string]string)
								// get service
								service, err := grpcClient.GetService(r.Params.Cluster, r.Params.Namespace, d.Spec.Host)
								if err != nil {
									log.Errorf("get service by grpc error: %v", err)
									return nil, err
								}
								// 获取service标签
								label = service.Spec.Selector
								deployment := &appsv1.DeploymentList{}
								// 把DR的标签添进去，就能找到对应的Deployment了
								for key, v := range dd.Labels {
									label[key] = v
								}
								// get deployment
								labelSelector := util.GenerateLabelSelector(label)
								deployment, err = grpcClient.GetDeployment(r.Params.Cluster, r.Params.Namespace, labelSelector)
								if err != nil {
									log.Errorf("get deployment by grpc error: %v", err)
								}
								for _, d := range deployment.Items {
									deployment := chart.Chart{}
									deployment.Name = fmt.Sprintf("部署: %s", d.Name)
									deployment.Rank = "deployment"
									destination.Children = append(destination.Children, deployment)
									// 计算Pod数量
									count++
								}
							}
						}
					}
				}
			}
			chartData.Children = append(chartData.Children, route)
		}
		chartData.Count = count
	}
	return &chartData, nil
}

func (r *VirtualServicesResource) GenerateCreateData(c *gin.Context) (err error) {
	switch r.Params.DataType {
	case "yaml":
		var j []byte
		create := common.PostType{}
		if err = c.BindJSON(&create); err != nil {
			return
		}
		if j, _, err = kit.YamlToJson(create.Context); err != nil {
			return
		}
		if err = json.Unmarshal(j, &r.PostData); err != nil {
			return
		}
	case "json":
		if err = c.BindJSON(&r.PostData); err != nil {
			return
		}
	default:
		return errors.New(common.ContentTypeError)
	}
	return nil
}
