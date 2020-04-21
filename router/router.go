package router

import (
	"github.com/gin-gonic/gin"
	"kingfisher/kf/common"
	jwtAuth "kingfisher/kf/middleware/jwt"
	"kingfisher/king-istio/impl"
	"net/http"
)

func SetupRouter(r *gin.Engine) *gin.Engine {

	//重新定义404
	r.NoRoute(NoRoute)

	authorize := r.Group("/", jwtAuth.JWTAuth())
	{
		// virtualServices
		authorize.GET(common.Istio+"virtualServices", impl.VirtualServices().List())
		authorize.GET(common.Istio+"virtualServicesChart/:name", impl.VirtualServices().Chart())
		authorize.GET(common.Istio+"virtualServices/:name", impl.VirtualServices().Get())
		authorize.DELETE(common.Istio+"virtualServices/:name", impl.VirtualServices().Delete())
		authorize.PATCH(common.Istio+"virtualServices/patch/:name", impl.VirtualServices().Patch())
		authorize.POST(common.Istio+"virtualServices", impl.VirtualServices().Create())
		authorize.PUT(common.Istio+"virtualServices", impl.VirtualServices().Update())
		// gateways
		authorize.GET(common.Istio+"gateways", impl.GateWays().List())
		authorize.GET(common.Istio+"gatewaysChart/:name", impl.GateWays().Chart())
		authorize.GET(common.Istio+"gateways/:name", impl.GateWays().Get())
		authorize.DELETE(common.Istio+"gateways/:name", impl.GateWays().Delete())
		authorize.PATCH(common.Istio+"gateways/patch/:name", impl.GateWays().Patch())
		authorize.POST(common.Istio+"gateways", impl.GateWays().Create())
		authorize.PUT(common.Istio+"gateways", impl.GateWays().Update())
		// destinationRules
		authorize.GET(common.Istio+"destinationRules", impl.DestinationRules().List())
		authorize.GET(common.Istio+"destinationRulesChart/:name", impl.DestinationRules().Chart())
		authorize.GET(common.Istio+"destinationRules/:name", impl.DestinationRules().Get())
		authorize.DELETE(common.Istio+"destinationRules/:name", impl.DestinationRules().Delete())
		authorize.PATCH(common.Istio+"destinationRules/patch/:name", impl.DestinationRules().Patch())
		authorize.POST(common.Istio+"destinationRules", impl.DestinationRules().Create())
		authorize.PUT(common.Istio+"destinationRules", impl.DestinationRules().Update())
		// envoyFilters
		authorize.GET(common.Istio+"envoyFilters", impl.EnvoyFilters().List())
		authorize.GET(common.Istio+"envoyFilters/:name", impl.EnvoyFilters().Get())
		authorize.DELETE(common.Istio+"envoyFilters/:name", impl.EnvoyFilters().Delete())
		authorize.PATCH(common.Istio+"envoyFilters/patch/:name", impl.EnvoyFilters().Patch())
		authorize.POST(common.Istio+"envoyFilters", impl.EnvoyFilters().Create())
		authorize.PUT(common.Istio+"envoyFilters", impl.EnvoyFilters().Update())
		// serviceEntries
		authorize.GET(common.Istio+"serviceEntries", impl.ServiceEntries().List())
		authorize.GET(common.Istio+"serviceEntries/:name", impl.ServiceEntries().Get())
		authorize.DELETE(common.Istio+"serviceEntries/:name", impl.ServiceEntries().Delete())
		authorize.PATCH(common.Istio+"serviceEntries/patch/:name", impl.ServiceEntries().Patch())
		authorize.POST(common.Istio+"serviceEntries", impl.ServiceEntries().Create())
		authorize.PUT(common.Istio+"serviceEntries", impl.ServiceEntries().Update())
		// sidecars
		authorize.GET(common.Istio+"sidecars", impl.Sidecars().List())
		authorize.GET(common.Istio+"sidecars/:name", impl.Sidecars().Get())
		authorize.DELETE(common.Istio+"sidecars/:name", impl.Sidecars().Delete())
		authorize.PATCH(common.Istio+"sidecars/patch/:name", impl.Sidecars().Patch())
		authorize.POST(common.Istio+"sidecars", impl.Sidecars().Create())
		authorize.PUT(common.Istio+"sidecars", impl.Sidecars().Update())

		// config
		authorize.GET(common.Istio+"rules", impl.Rules().List())
		authorize.GET(common.Istio+"rules/:name", impl.Rules().Get())
		authorize.DELETE(common.Istio+"rules/:name", impl.Rules().Delete())
		authorize.PATCH(common.Istio+"rules/patch/:name", impl.Rules().Patch())
		authorize.POST(common.Istio+"rules", impl.Rules().Create())
		authorize.PUT(common.Istio+"rules", impl.Rules().Update())
		// instance
		authorize.GET(common.Istio+"instances", impl.Instance().List())
		authorize.GET(common.Istio+"instances/:name", impl.Instance().Get())
		authorize.DELETE(common.Istio+"instances/:name", impl.Instance().Delete())
		authorize.PATCH(common.Istio+"instances/patch/:name", impl.Instance().Patch())
		authorize.POST(common.Istio+"instances", impl.Instance().Create())
		authorize.PUT(common.Istio+"instances", impl.Instance().Update())
	}
	return r
}

// 重新定义404错误
func NoRoute(c *gin.Context) {
	responseData := common.ResponseData{Code: http.StatusNotFound, Msg: "404 Not Found"}
	c.JSON(http.StatusNotFound, responseData)
}
