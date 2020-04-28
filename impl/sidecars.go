package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/open-kingfisher/king-istio/common/access"
	"github.com/open-kingfisher/king-istio/resource"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/handle"
)

type sidecars struct{}

func Sidecars() *sidecars {
	return &sidecars{}
}

func (v *sidecars) Get() func(c *gin.Context) {
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

func (v *sidecars) List() func(c *gin.Context) {
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

func (v *sidecars) Delete() func(c *gin.Context) {
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

func (v *sidecars) Patch() func(c *gin.Context) {
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

func (v *sidecars) Update() func(c *gin.Context) {
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

func (v *sidecars) Create() func(c *gin.Context) {
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

func (v *sidecars) newResource(c *gin.Context) (*resource.SidecarsResource, error) {
	// 获取clientSet，如果失败直接返回错误
	clientSet, err := access.IstioClient(c.Query("cluster"))
	if err != nil {
		return nil, err
	}
	// 获取HTTP的参数，存到handle.Resources结构体中
	commonParams := handle.GenerateCommonParams(c, nil)
	r := resource.SidecarsResource{
		Params: commonParams,
		Access: clientSet,
	}
	return &r, nil
}
