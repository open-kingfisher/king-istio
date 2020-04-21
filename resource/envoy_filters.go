package resource

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"kingfisher/kf/common"
	"kingfisher/kf/common/handle"
	"kingfisher/kf/common/log"
	"kingfisher/kf/kit"
	"kingfisher/king-istio/pkg/apis/networking/v1alpha3"
	networkingv1alpha3 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/networking/v1alpha3"
)

type EnvoyFiltersResource struct {
	Params   *handle.Resources
	PostData *v1alpha3.EnvoyFilter
	Access   *networkingv1alpha3.NetworkingV1alpha3Client
}

func (r *EnvoyFiltersResource) Get() (*v1alpha3.EnvoyFilter, error) {
	i, err := r.Access.EnvoyFilters(r.Params.Namespace).Get(r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "EnvoyFilter", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *EnvoyFiltersResource) List() (*v1alpha3.EnvoyFilterList, error) {
	return r.Access.EnvoyFilters(r.Params.Namespace).List(metav1.ListOptions{})
}

func (r *EnvoyFiltersResource) Delete() (err error) {
	if err = r.Access.VirtualServices(r.Params.Namespace).Delete(r.Params.Name, &metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.EnvoyFilters,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *EnvoyFiltersResource) Patch() (res *v1alpha3.EnvoyFilter, err error) {
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.EnvoyFilters(r.Params.Namespace).Patch(r.Params.Name, types.JSONPatchType, data); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.EnvoyFilters, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.EnvoyFilters,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *EnvoyFiltersResource) Update() (res *v1alpha3.EnvoyFilter, err error) {
	if res, err = r.Access.EnvoyFilters(r.Params.Namespace).Update(r.PostData); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.EnvoyFilters, err, r.PostData, r.PostData.Name)
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

func (r *EnvoyFiltersResource) Create() (res *v1alpha3.EnvoyFilter, err error) {
	if res, err = r.Access.EnvoyFilters(r.Params.Namespace).Create(r.PostData); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.EnvoyFilters, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.EnvoyFilters,
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

func (r *EnvoyFiltersResource) GenerateCreateData(c *gin.Context) (err error) {
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
