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
	"kingfisher/king-istio/pkg/apis/config/v1alpha2"
	configv1alpha2 "kingfisher/king-istio/pkg/client/clientset/versioned/typed/config/v1alpha2"
)

type RulesResource struct {
	Params   *handle.Resources
	PostData *v1alpha2.Rule
	Access   *configv1alpha2.ConfigV1alpha2Client
}

func (r *RulesResource) Get() (*v1alpha2.Rule, error) {
	i, err := r.Access.Rules(r.Params.Namespace).Get(r.Params.Name, metav1.GetOptions{})
	if err == nil {
		i.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "Rule", Version: "networking.istio.io/v1alpha3"})
	}
	return i, err
}

func (r *RulesResource) List() (*v1alpha2.RuleList, error) {
	return r.Access.Rules(r.Params.Namespace).List(metav1.ListOptions{})
}

func (r *RulesResource) Delete() (err error) {
	if err = r.Access.Rules(r.Params.Namespace).Delete(r.Params.Name, &metav1.DeleteOptions{}); err != nil {
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Rules,
		ActionType: common.Delete,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *RulesResource) Patch() (res *v1alpha2.Rule, err error) {
	var data []byte
	if data, err = json.Marshal(r.Params.PatchData.Patches); err != nil {
		return
	}
	if res, err = r.Access.Rules(r.Params.Namespace).Patch(r.Params.Name, types.JSONPatchType, data); err != nil {
		log.Errorf("%s patch error:%s; Json:%+v; Name:%s", common.Rules, err, string(data), r.Params.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Rules,
		ActionType: common.Patch,
		Resources:  r.Params,
		Name:       r.Params.Name,
	}
	if err = auditLog.InsertAuditLog(); err != nil {
		return
	}
	return
}

func (r *RulesResource) Update() (res *v1alpha2.Rule, err error) {
	if res, err = r.Access.Rules(r.Params.Namespace).Update(r.PostData); err != nil {
		log.Errorf("%s update error:%s; Json:%+v; Name:%s", common.Rules, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Rules,
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

func (r *RulesResource) Create() (res *v1alpha2.Rule, err error) {
	if res, err = r.Access.Rules(r.Params.Namespace).Create(r.PostData); err != nil {
		log.Errorf("%s create error:%s; Json:%+v; Name:%s", common.Rules, err, r.PostData, r.PostData.Name)
		return
	}
	auditLog := handle.AuditLog{
		Kind:       common.Rules,
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

func (r *RulesResource) GenerateCreateData(c *gin.Context) (err error) {
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
