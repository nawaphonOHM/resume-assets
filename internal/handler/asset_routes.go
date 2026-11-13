package handler

import "github.com/nawaphonOHM/whatever/v3/pkg/rest"

// RouteAssetPath is the single application route for asset retrieval.
const RouteAssetPath rest.Pathz = "/assets/:filename"

// NewAPIRegistration registers the single raw application route.
func NewAPIRegistration(handler *AssetHandler) *rest.RRestAPIRegistration {
	registration := &rest.RRestAPIRegistration{
		Version: 0,
		Prefix:  "",
		Group:   rest.APIGroupRaw,
	}
	if handler != nil {
		registration.Apis = []*rest.ExportableAPI{
			{
				Method:  rest.GET,
				Path:    RouteAssetPath,
				Handler: handler.ServeAsset,
			},
			{
				Method:  rest.HEAD,
				Path:    RouteAssetPath,
				Handler: handler.ServeAssetHead,
			},
		}
	}
	return registration
}

// RegisterRoutes creates the application blueprint with the asset route.
func RegisterRoutes(handler *AssetHandler) *rest.BluePrint {
	return rest.NewBluePrint().WithAPIs(NewAPIRegistration(handler))
}
