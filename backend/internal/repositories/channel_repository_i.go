package repositories

import "backend/internal/common"

type GetChannelsResult struct {
	Error    common.WrappedError
	Channels []string
}

type GetCurrentChannelResult struct {
	Error    common.WrappedError
	Channels string
}

type ChannelRepositoryI interface {
	AddChannel(iChannelName string, oDone chan common.WrappedError)
	GetChannels(oChannels chan GetChannelsResult)
	/// Does not check if channel exists
	SetCurrentChannel(channel string, oDone chan common.WrappedError)
	/// Returns NotFound if no current channel is set
	GetCurrentChannel(oDone chan GetCurrentChannelResult)
}
