package resource

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/kit"
	"istio.io/client-go/pkg/apis/config/v1alpha2"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type InstancesResource struct {
	Params   *handle.Resources
	PostData *v1alpha2.Instance
	Access   *versionedclient.Clientset
}

func (r *InstancesResource) Get() (*v1alpha2.Instance, error) {
	i, err := r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).Get(context.TODO(), r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "Rule", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *InstancesResource) List() (*v1alpha2.InstanceList, error) {
	return r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).List(context.TODO(), metav1.ListOptions{})
}

func (r *InstancesResource) Delete() (err error) {
	if err = r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).Delete(context.TODO(), r.Params.Name, metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Instances,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *InstancesResource) Patch() (res *v1alpha2.Instance, err error) {
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).Patch(context.TODO(), r.Params.Name, types.JSONPatchType, data, metav1.PatchOptions{}); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.Instances, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Instances,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *InstancesResource) Update() (res *v1alpha2.Instance, err error) {
	if res, err = r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).Update(context.TODO(), r.PostData, metav1.UpdateOptions{}); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.Instances, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Instances,
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

func (r *InstancesResource) Create() (res *v1alpha2.Instance, err error) {
	if res, err = r.Access.ConfigV1alpha2().Instances(r.Params.Namespace).Create(context.TODO(), r.PostData, metav1.CreateOptions{}); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.Instances, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Instances,
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

func (r *InstancesResource) GenerateCreateData(c *gin.Context) (err error) {
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
