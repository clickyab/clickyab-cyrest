package channel

import (
	"common/assert"
	"common/rabbit"
	"common/redis"
	"common/utils"
	"encoding/json"
	"modules/misc/middlewares"
	"modules/misc/trans"
	"modules/telegram/common/tgo"
	"modules/telegram/cyborg/commands"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"strconv"
	"time"

	"common/models/common"

	"modules/telegram/ad/ads"

	"fmt"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type channelPayload struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name" validate:"required"`
	Link   string `json:"link" `
}

// @Validate {
// }
type editChannelPayload struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link" `
}

//	createChannel
//	@Route	{
//	url	=	/
//	method	= post
//	payload	= channelPayload
//	resource = create_channel:self
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) createChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*channelPayload)
	m := ads.NewAdsManager()
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	if pl.UserID == 0 {
		pl.UserID = currentUser.ID
	}
	owner, err := aaa.NewAaaManager().FindUserByID(pl.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("create_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	ch := &ads.Channel{
		Name:          pl.Name,
		ArchiveStatus: ads.ArchiveStatusNo,
		AdminStatus:   ads.AdminStatusPending,
		Link:          common.MakeNullString(pl.Link),
		Active:        ads.ActiveStatusNo,
		UserID:        currentUser.ID,
	}
	assert.Nil(m.CreateChannel(ch))
	return u.OKResponse(ctx, ch)

}

//	getChannel
//	@Route	{
//	url	=	/:id
//	method	= get
//	resource = list_channel:self
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getChannel(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
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
//	payload	= editChannelPayload
//	resource = edit_channel:self
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) editChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*editChannelPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	channel, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	_, b := currentUser.HasPermOn("edit_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	channel.Name = pl.Name
	channel.Link = common.MakeNullString(pl.Link)

	m.UpdateChannel(channel)

	return u.OKResponse(ctx, channel)
}

//	active
//	@Route	{
//	url	=	/list/active_status/:id
//	method	= put
//	payload	= activePayload
//	resource = active_channel:global
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) active(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*activePayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	channel, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("active_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if pl.ActiveStatus == ads.ActiveStatusYes && channel.Active == ads.ActiveStatusNo {
		channel.Active = ads.ActiveStatusYes
	} else {
		channel.Active = ads.ActiveStatusNo
	}
	assert.Nil(m.UpdateChannel(channel))
	return u.OKResponse(ctx, channel)
}

//	changeArchive toggle archiving channel
//	@Route	{
//	url	=	/list/archive_status/:id
//	method	= put
//	payload	= archivePayload
//	resource = archive_channel:self
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeArchive(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*archivePayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	channel, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("archive_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if pl.ArchiveStatus == ads.ArchiveStatusYes && channel.ArchiveStatus == ads.ArchiveStatusNo {
		channel.ArchiveStatus = ads.ArchiveStatusYes
	} else {
		channel.ArchiveStatus = ads.ArchiveStatusNo
	}
	assert.Nil(m.UpdateChannel(channel))
	return u.OKResponse(ctx, channel)
}

// @Validate {
// }
type statusPayload struct {
	Status ads.AdminStatus `json:"status" validate:"required"`
}

// @Validate {
// }
type activePayload struct {
	ActiveStatus ads.ActiveStatus `json:"active_status" validate:"required"`
}

// @Validate {
// }
type archivePayload struct {
	ArchiveStatus ads.ArchiveStatus `json:"archive_status" validate:"required"`
}

// MsgInfo is msg info
type MsgInfo struct {
	CliID string
	Type  string
	Text  string
}

// GetLastResponse is the lst response command
type GetLastResponse struct {
	Data   []MsgInfo
	Status string
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

// Validate custom validation for user scope
func (lp *archivePayload) ValidateExtra(ctx echo.Context) error {
	if !lp.ArchiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("is invalid"),
		}
	}
	return nil
}

// Validate custom validation for user scope
func (lp *activePayload) ValidateExtra(ctx echo.Context) error {
	if !lp.ActiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("is invalid"),
		}
	}
	return nil
}

//	statusChannel the route for get status channel
//	@Route	{
//	url	=	/list/admin_status/:id
//	method	= put
//	payload	= statusPayload
//	resource = status_channel:parent
//	middleware = authz.Authenticate
//	200 = ads.Channel
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) statusChannel(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*statusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	cha, err := m.FindChannelByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)

	owner, err := aaa.NewAaaManager().FindUserByID(cha.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("status_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	cha.ID = id
	cha.AdminStatus = pl.Status
	assert.Nil(m.UpdateChannel(cha))
	return u.OKResponse(ctx, cha)

}

//	getLast get last messages for the specified channel
//	@Route	{
//	url	=	/last/:name/:count
//	method	= get
//	resource = get_last_channel:parent
//	middleware = authz.Authenticate
//	200 = GetLastResponse
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getLast(ctx echo.Context) error {
	var res []tgo.History
	var finalRes []MsgInfo

	count, err := strconv.Atoi(ctx.Param("count"))
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	//validate count between 1 - 99
	if count < 1 || count > 99 {
		return u.BadResponse(ctx, trans.E("count out of range"))
	}
	name := ctx.Param("name")
	hash := utils.Sha1(fmt.Sprintf("%s%d", name, count))
	//check if the key exists in redis
	b, err := aredis.GetHashKey(hash, "STATUS", true, 2*time.Hour)
	if b == "" || err != nil { //key not exists
		err = aredis.StoreHashKey(hash, "STATUS", "pending", 2*time.Hour)
		assert.Nil(err)
		rabbit.MustPublish(
			commands.GetLastCommand{
				Channel: name,
				Count:   count,
				HashKey: hash,
			},
		)
		return u.OKResponse(ctx, res)
	}
	if b == "pending" {
		return u.OKResponse(ctx, GetLastResponse{Status: "failed", Data: []MsgInfo{}})
	} else if b == "done" {
		stringRes, err := aredis.GetHashKey(hash, "DATA", true, 2*time.Hour)
		if err != nil {
			return u.BadResponse(ctx, trans.E("failed job"))
		}
		err = json.Unmarshal([]byte(stringRes), &res)
		if err != nil {
			return u.BadResponse(ctx, trans.E("failed job"))
		}
		for i := range res {
			if res[i].Media == nil {
				finalRes = append(finalRes, MsgInfo{CliID: res[i].ID, Text: res[i].Text, Type: tgo.Message})

			} else if res[i].Media.Type == tgo.Photo {
				finalRes = append(finalRes, MsgInfo{CliID: res[i].ID, Text: res[i].Media.Caption, Type: res[i].Media.Type})
			}
		}
		return u.OKResponse(ctx, GetLastResponse{
			Status: "done",
			Data:   finalRes,
		})
	} else if b == "failed" {
		aredis.RemoveKey(hash)
		return u.BadResponse(ctx, trans.E("failed job"))
	}
	return u.BadResponse(ctx, trans.E("failed job"))
}
