// Package v2 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package v2

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	BearerScopes = "Bearer.Scopes"
)

// Defines values for ImageStatusValue.
const (
	ImageStatusValueBuilding ImageStatusValue = "building"

	ImageStatusValueFailure ImageStatusValue = "failure"

	ImageStatusValuePending ImageStatusValue = "pending"

	ImageStatusValueRegistering ImageStatusValue = "registering"

	ImageStatusValueSuccess ImageStatusValue = "success"

	ImageStatusValueUploading ImageStatusValue = "uploading"
)

// Defines values for ImageTypes.
const (
	ImageTypesAws ImageTypes = "aws"

	ImageTypesAzure ImageTypes = "azure"

	ImageTypesEdgeCommit ImageTypes = "edge-commit"

	ImageTypesEdgeContainer ImageTypes = "edge-container"

	ImageTypesEdgeInstaller ImageTypes = "edge-installer"

	ImageTypesGcp ImageTypes = "gcp"

	ImageTypesGuestImage ImageTypes = "guest-image"

	ImageTypesImageInstaller ImageTypes = "image-installer"

	ImageTypesVsphere ImageTypes = "vsphere"
)

// Defines values for UploadStatusValue.
const (
	UploadStatusValueFailure UploadStatusValue = "failure"

	UploadStatusValuePending UploadStatusValue = "pending"

	UploadStatusValueRunning UploadStatusValue = "running"

	UploadStatusValueSuccess UploadStatusValue = "success"
)

// Defines values for UploadTypes.
const (
	UploadTypesAws UploadTypes = "aws"

	UploadTypesAwsS3 UploadTypes = "aws.s3"

	UploadTypesAzure UploadTypes = "azure"

	UploadTypesGcp UploadTypes = "gcp"
)

// AWSEC2UploadOptions defines model for AWSEC2UploadOptions.
type AWSEC2UploadOptions struct {
	Region            string   `json:"region"`
	ShareWithAccounts []string `json:"share_with_accounts"`
	SnapshotName      *string  `json:"snapshot_name,omitempty"`
}

// AWSEC2UploadStatus defines model for AWSEC2UploadStatus.
type AWSEC2UploadStatus struct {
	Ami    string `json:"ami"`
	Region string `json:"region"`
}

// AWSS3UploadOptions defines model for AWSS3UploadOptions.
type AWSS3UploadOptions struct {
	Region string `json:"region"`
}

// AWSS3UploadStatus defines model for AWSS3UploadStatus.
type AWSS3UploadStatus struct {
	Url string `json:"url"`
}

// AzureUploadOptions defines model for AzureUploadOptions.
type AzureUploadOptions struct {
	// Name of the uploaded image. It must be unique in the given resource group.
	// If name is omitted from the request, a random one based on a UUID is
	// generated.
	ImageName *string `json:"image_name,omitempty"`

	// Location where the image should be uploaded and registered.
	// How to list all locations:
	// https://docs.microsoft.com/en-us/cli/azure/account?view=azure-cli-latest#az_account_list_locations'
	Location string `json:"location"`

	// Name of the resource group where the image should be uploaded.
	ResourceGroup string `json:"resource_group"`

	// ID of subscription where the image should be uploaded.
	SubscriptionId string `json:"subscription_id"`

	// ID of the tenant where the image should be uploaded.
	// How to find it in the Azure Portal:
	// https://docs.microsoft.com/en-us/azure/active-directory/fundamentals/active-directory-how-to-find-tenant
	TenantId string `json:"tenant_id"`
}

// AzureUploadStatus defines model for AzureUploadStatus.
type AzureUploadStatus struct {
	ImageName string `json:"image_name"`
}

// ComposeId defines model for ComposeId.
type ComposeId struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectReference)
	ObjectReference `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
}

// ComposeMetadata defines model for ComposeMetadata.
type ComposeMetadata struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectReference)
	ObjectReference `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// ID (hash) of the built commit
	OstreeCommit *string `json:"ostree_commit,omitempty"`

	// Package list including NEVRA
	Packages *[]PackageMetadata `json:"packages,omitempty"`
}

// ComposeRequest defines model for ComposeRequest.
type ComposeRequest struct {
	Customizations *Customizations `json:"customizations,omitempty"`
	Distribution   string          `json:"distribution"`
	ImageRequest   ImageRequest    `json:"image_request"`
}

// ComposeStatus defines model for ComposeStatus.
type ComposeStatus struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectReference)
	ObjectReference `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	ImageStatus ImageStatus `json:"image_status"`
}

// Customizations defines model for Customizations.
type Customizations struct {
	Packages            *[]string     `json:"packages,omitempty"`
	PayloadRepositories *[]Repository `json:"payload_repositories,omitempty"`
	Subscription        *Subscription `json:"subscription,omitempty"`
	Users               *[]User       `json:"users,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectReference)
	ObjectReference `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
	Reason      string `json:"reason"`
}

// ErrorList defines model for ErrorList.
type ErrorList struct {
	// Embedded struct due to allOf(#/components/schemas/List)
	List `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Items []Error `json:"items"`
}

// GCPUploadOptions defines model for GCPUploadOptions.
type GCPUploadOptions struct {
	// Name of an existing STANDARD Storage class Bucket.
	Bucket string `json:"bucket"`

	// The name to use for the imported and shared Compute Engine image.
	// The image name must be unique within the GCP project, which is used
	// for the OS image upload and import. If not specified a random
	// 'composer-api-<uuid>' string is used as the image name.
	ImageName *string `json:"image_name,omitempty"`

	// The GCP region where the OS image will be imported to and shared from.
	// The value must be a valid GCP location. See https://cloud.google.com/storage/docs/locations.
	// If not specified, the multi-region location closest to the source
	// (source Storage Bucket location) is chosen automatically.
	Region string `json:"region"`

	// List of valid Google accounts to share the imported Compute Engine image with.
	// Each string must contain a specifier of the account type. Valid formats are:
	//   - 'user:{emailid}': An email address that represents a specific
	//     Google account. For example, 'alice@example.com'.
	//   - 'serviceAccount:{emailid}': An email address that represents a
	//     service account. For example, 'my-other-app@appspot.gserviceaccount.com'.
	//   - 'group:{emailid}': An email address that represents a Google group.
	//     For example, 'admins@example.com'.
	//   - 'domain:{domain}': The G Suite domain (primary) that represents all
	//     the users of that domain. For example, 'google.com' or 'example.com'.
	// If not specified, the imported Compute Engine image is not shared with any
	// account.
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
}

// GCPUploadStatus defines model for GCPUploadStatus.
type GCPUploadStatus struct {
	ImageName string `json:"image_name"`
	ProjectId string `json:"project_id"`
}

// ImageRequest defines model for ImageRequest.
type ImageRequest struct {
	Architecture  string        `json:"architecture"`
	ImageType     ImageTypes    `json:"image_type"`
	Ostree        *OSTree       `json:"ostree,omitempty"`
	Repositories  []Repository  `json:"repositories"`
	UploadOptions UploadOptions `json:"upload_options"`
}

// ImageStatus defines model for ImageStatus.
type ImageStatus struct {
	Status       ImageStatusValue `json:"status"`
	UploadStatus *UploadStatus    `json:"upload_status,omitempty"`
}

// ImageStatusValue defines model for ImageStatusValue.
type ImageStatusValue string

// ImageTypes defines model for ImageTypes.
type ImageTypes string

// List defines model for List.
type List struct {
	Kind  string `json:"kind"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
}

// OSTree defines model for OSTree.
type OSTree struct {
	Ref *string `json:"ref,omitempty"`
	Url *string `json:"url,omitempty"`
}

// ObjectReference defines model for ObjectReference.
type ObjectReference struct {
	Href string `json:"href"`
	Id   string `json:"id"`
	Kind string `json:"kind"`
}

// PackageMetadata defines model for PackageMetadata.
type PackageMetadata struct {
	Arch      string  `json:"arch"`
	Epoch     *string `json:"epoch,omitempty"`
	Name      string  `json:"name"`
	Release   string  `json:"release"`
	Sigmd5    string  `json:"sigmd5"`
	Signature *string `json:"signature,omitempty"`
	Type      string  `json:"type"`
	Version   string  `json:"version"`
}

// Repository defines model for Repository.
type Repository struct {
	Baseurl    *string `json:"baseurl,omitempty"`
	CheckGpg   *bool   `json:"check_gpg,omitempty"`
	GpgKey     *string `json:"gpg_key,omitempty"`
	IgnoreSsl  *bool   `json:"ignore_ssl,omitempty"`
	Metalink   *string `json:"metalink,omitempty"`
	Mirrorlist *string `json:"mirrorlist,omitempty"`
	Rhsm       bool    `json:"rhsm"`
}

// Subscription defines model for Subscription.
type Subscription struct {
	ActivationKey string `json:"activation_key"`
	BaseUrl       string `json:"base_url"`
	Insights      bool   `json:"insights"`
	Organization  string `json:"organization"`
	ServerUrl     string `json:"server_url"`
}

// UploadOptions defines model for UploadOptions.
type UploadOptions interface{}

// UploadStatus defines model for UploadStatus.
type UploadStatus struct {
	Options interface{}       `json:"options"`
	Status  UploadStatusValue `json:"status"`
	Type    UploadTypes       `json:"type"`
}

// UploadStatusValue defines model for UploadStatusValue.
type UploadStatusValue string

// UploadTypes defines model for UploadTypes.
type UploadTypes string

// User defines model for User.
type User struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectReference)
	ObjectReference `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Groups *[]string `json:"groups,omitempty"`
	Key    *string   `json:"key,omitempty"`
	Name   string    `json:"name"`
}

// Page defines model for page.
type Page string

// Size defines model for size.
type Size string

// PostComposeJSONBody defines parameters for PostCompose.
type PostComposeJSONBody ComposeRequest

// GetErrorListParams defines parameters for GetErrorList.
type GetErrorListParams struct {
	// Page index
	Page *Page `json:"page,omitempty"`

	// Number of items in each page
	Size *Size `json:"size,omitempty"`
}

// PostComposeJSONRequestBody defines body for PostCompose for application/json ContentType.
type PostComposeJSONRequestBody PostComposeJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create compose
	// (POST /compose)
	PostCompose(ctx echo.Context) error
	// The status of a compose
	// (GET /composes/{id})
	GetComposeStatus(ctx echo.Context, id string) error
	// Get the metadata for a compose.
	// (GET /composes/{id}/metadata)
	GetComposeMetadata(ctx echo.Context, id string) error
	// Get a list of all possible errors
	// (GET /errors)
	GetErrorList(ctx echo.Context, params GetErrorListParams) error
	// Get error description
	// (GET /errors/{id})
	GetError(ctx echo.Context, id string) error
	// Get the openapi spec in json format
	// (GET /openapi)
	GetOpenapi(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostCompose converts echo context to params.
func (w *ServerInterfaceWrapper) PostCompose(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostCompose(ctx)
	return err
}

// GetComposeStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposeStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposeStatus(ctx, id)
	return err
}

// GetComposeMetadata converts echo context to params.
func (w *ServerInterfaceWrapper) GetComposeMetadata(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetComposeMetadata(ctx, id)
	return err
}

// GetErrorList converts echo context to params.
func (w *ServerInterfaceWrapper) GetErrorList(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetErrorListParams
	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "size" -------------

	err = runtime.BindQueryParameter("form", true, false, "size", ctx.QueryParams(), &params.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter size: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetErrorList(ctx, params)
	return err
}

// GetError converts echo context to params.
func (w *ServerInterfaceWrapper) GetError(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetError(ctx, id)
	return err
}

// GetOpenapi converts echo context to params.
func (w *ServerInterfaceWrapper) GetOpenapi(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetOpenapi(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/compose", wrapper.PostCompose)
	router.GET(baseURL+"/composes/:id", wrapper.GetComposeStatus)
	router.GET(baseURL+"/composes/:id/metadata", wrapper.GetComposeMetadata)
	router.GET(baseURL+"/errors", wrapper.GetErrorList)
	router.GET(baseURL+"/errors/:id", wrapper.GetError)
	router.GET(baseURL+"/openapi", wrapper.GetOpenapi)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xbe28bOZL/KkTvAZ7BdUuyHo4jYDDrON6s9yYPWM4s7mJDoNglNTfdZIdkW1ECf/cD",
	"H93qB2XJO55ZDOB/Ykkkq35VrCpWFZnvAeFZzhkwJYPp9yDHAmegQLhvK9B/Y5BE0FxRzoJp8AGvAFEW",
	"w9cgDOArzvIUGtPvcFpAMA2Og/v7MKB6zZcCxCYIA4YzPWJmhoEkCWRYL1GbXP8ulaBsZZZJ+s3D+12R",
	"LUAgvkRUQSYRZQgwSZAjWEdTEqjQDAY78Zi5D+G5LwcN6bN/zi7Ohx/zlOP4vYFm5Rc8B6Go5S9gZTB/",
	"L1EF0wCKaA1SRcdB2GYRBjLBAuZrqpI5JoQXbkuq1Z+C4+FoPDl5cfpycDwMbsPA6MADtyKOhcAbQ5vh",
	"XCZcza3AdUzZJipHu6juw0DAl4IKiDUAJ5Mf6221mi/+BURpvnVNzRRWhUdROKNNRDij0YCcjgYvXo5e",
	"vJhMXk7i8cKnsUequCWM5lvR2AF+NnraXfbrcw/zXYorROr3nToLPclL/1shYI9wNMMrqEym5Yk4A+2H",
	"KgFUGDIQI7Oghy4Vygqp0AJQweiXQocLM3FF74AhAZIXggBaCV7kvRt2uUSaCaIS8YwqBTFaCp6ZJVoW",
	"kCpEGAnMYp4hzgAtsIQYcYYw+vjx8jWi8oatgIHACuLeDdvGAmvhBpjPhFJOsHI72BTwFzeC1gkIMFgM",
	"FSQTXqSxEa6UG7MY6b2UCoTh/3e+RoqjlEqFcJqiko2c3rBEqVxO+/2YE9nLKBFc8qXqEZ71gUWF7JOU",
	"9rHenr7zrZ/vKKx/Mj9FJKVRihVI9Rf8rXS+uWY0r5gctRSgrREKvbV+L7LbMTfb8fBON7fuANW09+Ka",
	"FwSzK0fmjeHoi4XFooIwp3EX1OVrDak+7d8AM4ZJfLoYkggvhuNoPD4eRS8HZBKdHA9HgxM4HbyEoQ+d",
	"AoaZegCXBmEnHYbKmcuSshhRVXqLcVH0gQuF00PsprQZRe8giqkAorjY9JcFi3EGTOFUdkajhK8jxSPN",
	"OrKQW0qakBewnCxOomMyWkbjGA8ifDIcRoPF4GQwHL2MX8Qv9ga6rca6e9uxwJpX7olcuyJjM3AdEgla",
	"eGsEfBDOddIk4dIYAE7T98tg+ul78F8ClsE0+Et/m1T1XdrQf28WX8ESBDACwX3YAR03wR4PR6CP+whO",
	"Xy6i42E8ivB4chKNhycnk8l4PBgMBkEYLLnIsAqmQVEYZe4RLPYIdLsV6S0oHGOFn1IwLpUAmBOeZVR5",
	"XeaHBMvkx9JzFgVNFXLTPe6XY/IZryztdmpqRmzcpYykRUzZCr27+PXqLKjlSw/J42hUiuhkU/cP6e/K",
	"HlddkySFVDyj33B11j4E4rw5+z4MYqoVsChUJ90QCaTRqU9R1orFFtJDLC/15BJ+22wa3NuEa+JvHfLJ",
	"vMKwkhXdvSI4CH6PdnR2uEBni5pQ6nZXy8pzLtVKgHxcRp7jjY5gcwE5l1RxUcp7iI1elYs23mS/FmD3",
	"UZrV596HQSFd7XcQjo8SxCEOEgYXQnDxlHZBeAxeRetJuJY3ePIdLK1iHg6VhkM1vUXYb0FGyl+odbbD",
	"JDWzPWZfqv+gfbDa9W1EwwUMKT/yN+cf9hQDi4J8BrU7PcQMwVcqlQ64s+uzd6/Prl6jmeJCB2SSYinR",
	"K0Oi107O3ZfIcdgZyPyFyHUCtnpQHBUS0JILl27lXCiXnJt6NUY6ShUK0AVbUeYyst4Nu66yM0OoVbvo",
	"KtdlZG/OP6BccK22EK0TShJdsxQS4htW8n0/c7RsfmfYWyw9pAsdrpDMgdAl1dhcUXPDjoiNoCLCOY1u",
	"isFgRPSJbj7BEbLKKNkhLGs5pUb9mKJnW7R2ValFtOO11LWSaU3TVKumUq7idf3qqs3p07RdKlVi/Z3G",
	"hnqZ3PXQDACVWS1JeRH3VpyvUjA5rbSmY9LdflXauGqxrsTQQMyKVNHIIS+nI5JyCVJpmHqSTTNv2A+u",
	"iinN0xpmtexHrWaScAkM4ULxDCtKcJpu2kqG4hGNnFZ5qVMUviz1YuRG5XSN11BpWrLPfI159m7YBSZJ",
	"aSRG64QzhamukEtNiTLBcmyQRt5DvxoENo2UCAuY3jCEInSkz4Lpd8gwTWl8fzRFZwyZbwjHsQCpTRAr",
	"JCAXIHU82vIimgRqidVDf+MCOe2F6AinlMBf3Xe950c9x1mCuKMEzuy6R2KwrB2JXbyzTcRVYrwt/yvO",
	"c5lz1Vu5ReWaOiRTmjxWG07+ss+hcbVUEGeUSa8OYp5hyqbf7V/N0LgnmhVUAbK/oh9yQTMsNj92maep",
	"ZWgaNPpUt7uPlVvb1sjW9Y4QF+iohcnvdQ+bJpV2jQ0O2lARZpsbVuq36U2fTPIx7ViFLhmb9nDo5gVh",
	"YLetq+YgDJyC6z8+IoPb1Rl1h5ivaqzO2KcrW8PAHUfzdvWIJQEWY6aihcA0jkaD0eR4tLc+rJEL91XB",
	"jYqh29YVJKEKiCpES5yvpyfzk/Huc97+fECqf73JwRRHtsLct+b97FrPMhI/edJtT/s5zw+q75q5Vqcz",
	"XVddQyst6B22t+W27DKxRxdSv5rrk62AhxFo2HlbvLIIa2K1jLShsCIz0wpCQGohl5imVhU5MF3RGz+j",
	"qftokdnPZRdWf7v1WFjNbmqs8FqzMf0zHZHiFURV+8F9M4cpiPIHyqTCaWp+WJFc/6vdoPJT87cx607m",
	"Op/yoipLhuZefabMX8GUF21ugDIFK1uIlZde3RHFFU59Q63NMUzD6obOXozZxeHOCiIMnG957keW3W5F",
	"/7RvY0Bf69IXCHZebXQZtyrFDoLEQegGG79yd2i920cLS10ZDj6ltFtJ3hjpBQE53zFSng6epD4FLP1j",
	"kq6yeLJriOEyRu848zwDdyAkPaSKdmHLwN4u28INrRIqjDoq1CJttwzFEpx1bI2qKiJi1hMQJ9i2xbXX",
	"AlP9mErV14Z3urU8TYfLPpf9Rg9VpD5zJAmQz/NVvqrJu+A8BWyaJqt8Nf8MG7+VrRgXMJcy9a/NQOGU",
	"ss9+gTKqK3vZW0LMBXaHc4+LVb9c97M+EH6y49FoqMvF4YlW6U/VMbtPOsskdTGoCaLCoId7BJji0vD/",
	"2W3gT6eRPnxxVuOM9b8nY/uLwfcKS3g/OwCLSGTmU1Q73dLTfC43a/W+Wv5GFL2zPRy3X817cCACVKSH",
	"akhzLOWai9gHV1vR3GuOXWs8QHrKJF0lrXt/JQoIPZbDxQoz161s8h8OxoPR0Jth6SQZRBdyvWfY09qt",
	"Id+bNDaQhG0tN5jWVFYT17eTnXYUZ3BAP833NuM+3LumfdG/b0mnX7aXR/e+3TTeHq4I+G8Rv0y/Dpf+",
	"wBXtQuYRspcrtOiPTyWrZPSQEsEudDWCPwUNy+Opnj93GR6clIqCsV2ZZx1ON/Vcy54cVbmkzUS9VCQ8",
	"aRvd1MfNGmgbFMyg961Su/rpRFMpkwji4WRy/BKdnZ2dnY/efcPnx+n/vb48fnd9MdG/Xb4Tb/7nQrz9",
	"X/rfb99+XBd/x1dn/8iufuGX366Wwy+vh/HrybfBq+uv/ZOvPhDdQrmQIPa/utlR0N6aV15ACkHVZqY1",
	"aFX0CrCwSl+YT38rg/g//nldPhozodnOq+jqU8A+HaNsybsdwJnrUClubjxdp9hWDLaBIntBGKSUALN5",
	"nXutdpZjkgAa9gaBy5SrfGG9XvewGTaHtFsr+79cnl+8m11Ew96gl6gsNXtIlVHa+9krw95d4QlkWrEI",
	"57SWsE2DobtcYXpgGox6g96xKRRUYtTUdw1sE8S49NwUnAvAChBGDNbIzQ5RznWORnGabhDhTLorBL5E",
	"Eu5A4FIXRj2up27e/NmeLhUoBr3E9YfrFzWXcTANPnCpnGiBtQOQ6hWPN/YWyWSIxqPyPKW2/9v/l7sg",
	"2j4IfPCytnn1e9+0N31821c2Odd7oakNB8dPzf0ytoxbKreDKMESSYWFglhv43gweDL+7u6py/uS2d62",
	"2+nyJZflf/z78z8rlDaSz8AQlYhaNJb76Pfn/pHhQiVc0G/2liQHobM/VBmnRTL+I5B8ZnzNqn2wSpj8",
	"ESbwkcHXHIiCGIGegzghhdBuUY+15hgro+yn2/vbMJBFlmFd/pVBowwuel0ZaWT/O43vzSnmu5h8A8pe",
	"+piT3FxRIndAIy4MxRQ0NEfOXFwZSyFpEYNE6wRUAkJPZtzSKnVo0gCIIe7Gmzegmo8hwsar6k/+F2MV",
	"YQtWcbQyV6HmtbKOsdvHyu7JVD2+1J8uP/kDottO8Bo8dfCqGoUdC2rq5T8Wu8rA8Ry2nsPWQWHruhV4",
	"dscv08kp24MPBrJyoqW4pIzKpBW+AMFXTBTSGaf2asoZEqAKwSBGMehKRSLO6i+ry2fb9jb4gXBWtTGf",
	"A9regLZ9Pdi1ruv6VpavRuzL+HIrn+Pcc5z7c8S5TmzSBo1rhqzjnSEua/GtE2K2D+c6wcUn2XZK31xU",
	"7WpA1eaZm6zf1fW3Mvis3b5J5kvklPHsZv8ZN7OG/udzMlwZEE5TlHMp6SKFypq2bra/KMLMtpkYqf5f",
	"j0W2fZe42CBzdPod9bAMoKL7W0/90R98hldb+eyjzz76GB+1a+ukjV9WTdPd5997N8Vv1U2wjpzxVkQZ",
	"0jpwzzf/jJnDg+LcV1eWvjjz1j2B5HFB7LtdO7fTFsc57Wk+MqHuf8zhnPbtGx3TewcRle+v+3dDk0+0",
	"mvUKryhbPcRAKryC38jGKJGVTzQrNvvo3N7/fwAAAP//5y+av8k/AAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
