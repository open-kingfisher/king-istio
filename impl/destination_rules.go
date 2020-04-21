package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"kingfisher/kf/common"
	accessK8s "kingfisher/kf/common/access"
	"kingfisher/kf/common/handle"
	"kingfisher/king-istio/common/access"
	"kingfisher/king-istio/resource"
)

type destinationRules struct{}

func DestinationRules() *destinationRules {
	return &destinationRules{}
}

func (v *destinationRules) Get() func(c *gin.Context) {
	return func(c *gin.Context) {
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		response, err := r.Get()
		responseData := handle.HandlerResponse(response, err)
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) List() func(c *gin.Context) {
	return func(c *gin.Context) {
		var responseData *common.ResponseData
		r, err := v.newResource(c)
		if err != nil {
			responseData = handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData.Data)
			return
		}
		response, err := r.List()
		if err != nil {
			responseData = handle.HandlerResponse(nil, err)
		} else {
			responseData = handle.HandlerResponse(response.Items, err)
		}
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		err = r.Delete()
		responseData := handle.HandlerResponse(nil, err)
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) Patch() func(c *gin.Context) {
	return func(c *gin.Context) {
		var responseData *common.ResponseData
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		if err := c.BindJSON(&r.Params.PatchData); err == nil {
			response, err := r.Patch()
			responseData = handle.HandlerResponse(response, err)
		} else {
			responseData = handle.HandlerResponse(nil, err)
		}
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) Update() func(c *gin.Context) {
	return func(c *gin.Context) {
		var responseData *common.ResponseData
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		if err := c.BindJSON(&r.PostData); err == nil {
			response, err := r.Update()
			responseData = handle.HandlerResponse(response, err)
		} else {
			responseData = handle.HandlerResponse(nil, err)
		}
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) Create() func(c *gin.Context) {
	return func(c *gin.Context) {
		var responseData *common.ResponseData
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		if err := r.GenerateCreateData(c); err == nil {
			if r.PostData != nil {
				response, err := r.Create()
				responseData = handle.HandlerResponse(response, err)
			} else {
				responseData = handle.HandlerResponse(nil, errors.New("the post data does not match the type"))
			}
		} else {
			responseData = handle.HandlerResponse(nil, err)
		}
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) Chart() func(c *gin.Context) {
	return func(c *gin.Context) {
		r, err := v.newResource(c)
		if err != nil {
			responseData := handle.HandlerResponse(nil, err)
			c.JSON(responseData.Code, responseData)
			return
		}
		response, err := r.Chart()
		responseData := handle.HandlerResponse(response, err)
		c.JSON(responseData.Code, responseData)
	}
}

func (v *destinationRules) newResource(c *gin.Context) (*resource.DestinationRulesResource, error) {
	// 获取clientSet，如果失败直接返回错误
	clientSet, err := access.IstioNetworkingClient(c.Query("cluster"))
	if err != nil {
		return nil, err
	}
	// 获取HTTP的参数，存到handle.Resources结构体中
	// 获取clientSet，如果失败直接返回错误
	clientSetK8s, err := accessK8s.Access(c.Query("cluster"))
	if err != nil && err.Error() == common.ClusterNotExistError {
		err = errors.New("cluster does not exist")
	}
	commonParams := handle.GenerateCommonParams(c, clientSetK8s)
	r := resource.DestinationRulesResource{
		Params: commonParams,
		Access: clientSet,
	}
	return &r, nil
}
