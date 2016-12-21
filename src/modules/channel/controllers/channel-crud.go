package channel

import (
	"modules/channel/chn"

	"modules/user/middlewares"

	"modules/user/aaa"

	"modules/misc/trans"
	"net/http"

	"strconv"

	"common/assert"

	"modules/misc/middlewares"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type channelPayload struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name" validate:"required"`
	Link   string `json:"link" `
	Admin  string `json:"admin"`
}

//	createChannel
//	@Route	{
//	url	=	/create
//	method	= post
//	payload	= channelPayload
//	resource = create_channel:self
//	middleware = authz.Authenticate
//	200 = chn.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) createChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*channelPayload)
	m := chn.NewChnManager()
	if pl.UserID == 0 {
		usr, _ := authz.GetUser(ctx)
		pl.UserID = usr.ID
	}
	user, err := aaa.NewAaaManager().FindUserByID(pl.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := user.HasPermOn("create_channel", pl.UserID, user.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	ch := m.Create(pl.Admin, pl.Link, pl.Name, chn.ChannelStatusPending, pl.UserID)
	return u.OKResponse(ctx, ch)

}

//	getChannel
//	@Route	{
//	url	=	/:id
//	method	= get
//	resource = list_channel:self
//	middleware = authz.Authenticate
//	200 = chn.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getChannel(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := chn.NewChnManager()
	channel, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("list_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	return u.OKResponse(ctx, channel)
}

//	editChannel
//	@Route	{
//	url	=	/:id
//	method	= put
//	payload	= channelPayload
//	resource = edit_channel:self
//	middleware = authz.Authenticate
//	200 = chn.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) editChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*channelPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := chn.NewChnManager()
	channel, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("edit_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	ch := m.EditChannel(pl.Admin, pl.Link, pl.Name, channel.Status, owner.ID, id)

	return u.OKResponse(ctx, ch)
}

// @Validate {
// }
type statusPayload struct {
	Status chn.ChannelStatus `json:"status" validate:"required"`
}

// Validate custom validation for user scope
func (lp *statusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.Status.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("is invalid"),
		}
	}
	return nil
}

//	StatusChannel
//	@Route	{
//	url	=	/status/:id
//	method	= put
//	payload	= statusPayload
//	resource = status_channel:parent
//	middleware = authz.Authenticate
//	200 = chn.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) StatusChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*statusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := chn.NewChnManager()
	cha, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}

	owner, err := aaa.NewAaaManager().FindUserByID(cha.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("status_channel", owner.ID, owner.ParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	cha.ID = id
	cha.Status = pl.Status
	assert.Nil(m.UpdateChannel(cha))
	return u.OKResponse(ctx, cha)

}
