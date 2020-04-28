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
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type SidecarsResource struct {
	Params   *handle.Resources
	PostData *v1alpha3.Sidecar
	Access   *versionedclient.Clientset
}

func (r *SidecarsResource) Get() (*v1alpha3.Sidecar, error) {
	var ctx context.Context
	i, err := r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).Get(ctx, r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "Sidecar", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *SidecarsResource) List() (*v1alpha3.SidecarList, error) {
	var ctx context.Context
	return r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).List(ctx, metav1.ListOptions{})
}

func (r *SidecarsResource) Delete() (err error) {
	var ctx context.Context
	if err = r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).Delete(ctx, r.Params.Name, metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Sidecars,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *SidecarsResource) Patch() (res *v1alpha3.Sidecar, err error) {
	var ctx context.Context
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).Patch(ctx, r.Params.Name, types.JSONPatchType, data, metav1.PatchOptions{}); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.Sidecars, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Sidecars,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *SidecarsResource) Update() (res *v1alpha3.Sidecar, err error) {
	var ctx context.Context
	if res, err = r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).Update(ctx, r.PostData, metav1.UpdateOptions{}); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.Sidecars, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Sidecars,
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

func (r *SidecarsResource) Create() (res *v1alpha3.Sidecar, err error) {
	var ctx context.Context
	if res, err = r.Access.NetworkingV1alpha3().Sidecars(r.Params.Namespace).Create(ctx, r.PostData, metav1.CreateOptions{}); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.Sidecars, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Sidecars,
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

func (r *SidecarsResource) GenerateCreateData(c *gin.Context) (err error) {
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
