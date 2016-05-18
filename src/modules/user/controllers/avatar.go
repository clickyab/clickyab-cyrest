package user

import (
	"bytes"
	"common/assert"
	"image/png"
	"strconv"

	"modules/user/assets"

	"modules/user/aaa"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

// getAvatar get the avatar from the system. default is the gravatar image
// @Route {
// 		url = /avatar/:user_id/:size/avatar.png
//		method = get
//      #payload = payload class
//		#resource = resource_name
//		produce = image/png
//		_size_ = integer, size of the image
//		_user_id_ = integer, the user id to get the avatar
//      200 = base.NormalResponse
//      404 = base.ErrorResponseSimple
// }
func (u *Controller) getAvatar(ctx *gin.Context) {
	UserID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 0)
	if err != nil {
		u.NotFoundResponse(ctx, nil)
		return
	}
	_, err = aaa.NewAaaManager().FindUserByID(UserID)
	if err != nil {
		u.NotFoundResponse(ctx, nil)
		return
	}

	size, err := strconv.ParseInt(ctx.Param("size"), 10, 0)
	if err != nil {
		size = 512
	}
	if size != 16 && size != 32 && size != 64 && size != 128 && size != 256 && size != 512 {
		u.NotFoundResponse(ctx, nil)
		return
	}
	data, err := assets.Asset("data/default.png")
	assert.Nil(err)

	buffer := bytes.NewBuffer(data)
	img, err := png.Decode(buffer)
	assert.Nil(err)

	dst := imaging.Resize(img, int(size), int(size), imaging.Linear)
	buffer.Reset()

	_ = png.Encode(buffer, dst)
	ctx.Header("Content-Type", "image/png")
	_, _ = ctx.Writer.Write(buffer.Bytes())
}
