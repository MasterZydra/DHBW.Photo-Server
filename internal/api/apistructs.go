package api

import (
	"DHBW.Photo-Server/internal/image"
	"DHBW.Photo-Server/internal/order"
	"time"
)

// These structs are used for easier consumption of the backend API.
// For each POST request there is a RequestData object which defines the parameters of this request.
// For each request there is a ResultData object which always has an Error field and some other result fields
// to send back to the API client.
// The Error field contains any error from the server that could have happened during the API call.

type BaseRes interface {
	GetError() string
}

type TestReqData struct {
	SomeString string
	SomeInt    int
}
type TestResData struct {
	Error            string
	SomeResultString string
	SomeResultInt    int
}

func (a TestResData) GetError() string {
	return a.Error
}

type RegisterReqData struct {
	Username             string
	Password             string
	PasswordConfirmation string
}
type RegisterResData struct {
	Error string
}

func (a RegisterResData) GetError() string {
	return a.Error
}

type UploadReqData struct {
	Base64Image  string
	Filename     string
	CreationDate time.Time
}
type UploadResData struct {
	Error string
}

func (a UploadResData) GetError() string {
	return a.Error
}

type ImageResData struct {
	Error string
	Image *image.Image
}

func (a ImageResData) GetError() string {
	return a.Error
}

type ThumbnailsReqData struct {
	Index  int
	Length int
}
type ThumbnailsResData struct {
	Error       string
	Images      []*image.Image
	TotalImages int
}

func (a ThumbnailsResData) GetError() string {
	return a.Error
}

type OrderListResData struct {
	OrderList []*order.ListEntry
	Error     string
}

func (a OrderListResData) GetError() string {
	return a.Error
}

type AddOrderListEntryReqData struct {
	ImageName string
}
type AddOrderListEntryResData struct {
	Error string
}

func (a AddOrderListEntryResData) GetError() string {
	return a.Error
}

type RemoveOrderListEntryReqData struct {
	ImageName string
}
type RemoveOrderListEntryResData struct {
	Error string
}

func (a RemoveOrderListEntryResData) GetError() string {
	return a.Error
}

type ChangeOrderListEntryReqData struct {
	ImageName      string
	Format         string
	NumberOfPrints int
}
type ChangeOrderListEntryResData struct {
	Error string
}

func (a ChangeOrderListEntryResData) GetError() string {
	return a.Error
}

type DeleteOrderListReqData struct {
}
type DeleteOrderListResData struct {
	Error string
}

func (a DeleteOrderListResData) GetError() string {
	return a.Error
}

type DownloadOrderListResData struct {
	Error         string
	Base64ZipFile string
}

func (a DownloadOrderListResData) GetError() string {
	return a.Error
}
