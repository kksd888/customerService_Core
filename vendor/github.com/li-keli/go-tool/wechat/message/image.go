package message

//Image 图片消息
type Image struct {
	CommonToken

	Image struct {
		MediaID string `json:"media_id" xml:"MediaId"`
	} `json:"image" xml:"Image"`
}

//NewImage 回复图片消息
func NewImage(mediaID string) *Image {
	image := new(Image)
	image.Image.MediaID = mediaID
	return image
}
