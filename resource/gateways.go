package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"kingfisher/kf/common"
	"kingfisher/kf/common/handle"
	"kingfisher/kf/common/log"
	"kingfisher/kf/kit"
	"kingfisher/king-istio/common/chart"
	"kingfisher/king-istio/pkg/apis/networking/v1alpha3"
	networkingv1alpha3 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/networking/v1alpha3"
)

type GatewaysResource struct {
	Params   *handle.Resources
	PostData *v1alpha3.Gateway
	Access   *networkingv1alpha3.NetworkingV1alpha3Client
}

func (r *GatewaysResource) Get() (*v1alpha3.Gateway, error) {
	i, err := r.Access.Gateways(r.Params.Namespace).Get(r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "Gateway", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *GatewaysResource) List() (*v1alpha3.GatewayList, error) {
	return r.Access.Gateways(r.Params.Namespace).List(metav1.ListOptions{})
}

func (r *GatewaysResource) Delete() (err error) {
	if err = r.Access.Gateways(r.Params.Namespace).Delete(r.Params.Name, &metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Gateways,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *GatewaysResource) Patch() (res *v1alpha3.Gateway, err error) {
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.Gateways(r.Params.Namespace).Patch(r.Params.Name, types.JSONPatchType, data); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.Gateways, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Gateways,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *GatewaysResource) Update() (res *v1alpha3.Gateway, err error) {
	if r.Params.PostType == "form" {
		if gateway, err := r.Access.Gateways(r.Params.Namespace).Get(r.PostData.Name, metav1.GetOptions{}); err != nil {
			log.Errorf("%s get error:%s; Json:%+v; Name:%s", common.Gateways, err, r.PostData, r.PostData.Name)
			return nil, err
		} else {
			for i, new := range r.PostData.Spec.Servers {
				for _, g := range gateway.Spec.Servers {
					if new.Port.Name == g.Port.Name {
						r.PostData.Spec.Servers[i] = g
						r.PostData.Spec.Servers[i].Port.Protocol = new.Port.Protocol
						r.PostData.Spec.Servers[i].Port.Number = new.Port.Number
						r.PostData.Spec.Servers[i].Hosts = new.Hosts
					}
				}
			}
		}
	}
	if res, err = r.Access.Gateways(r.Params.Namespace).Update(r.PostData); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.Gateways, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Gateways,
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

func (r *GatewaysResource) Create() (res *v1alpha3.Gateway, err error) {
	if res, err = r.Access.Gateways(r.Params.Namespace).Create(r.PostData); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.Gateways, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Gateways,
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

func (r *GatewaysResource) Chart() (interface{}, error) {
	chartData := chart.Chart{}
	if gateway, err := r.Get(); err != nil {
		return nil, err
	} else {
		count := 0
		chartData.Name = r.Params.Name
		chartData.Rank = "gateway"
		for _, v := range gateway.Spec.Servers {
			port := chart.Chart{}
			port.Name = fmt.Sprintf("名称:%s 端口:%d 协议:%s", v.Port.Name, v.Port.Number, v.Port.Protocol)
			port.Rank = "port"
			for _, p := range v.Hosts {
				hosts := chart.Chart{}
				hosts.Name = p
				hosts.Rank = "hosts"
				// 计算Pod数量
				count++
				port.Children = append(port.Children, hosts)
			}
			chartData.Children = append(chartData.Children, port)
		}
		chartData.Count = count
	}
	return &chartData, nil
}

func (r *GatewaysResource) GenerateCreateData(c *gin.Context) (err error) {
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
