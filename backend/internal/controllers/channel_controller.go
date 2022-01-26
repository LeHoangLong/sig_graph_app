package controllers

import "backend/internal/services"

type ChannelController struct {
	channelService *services.ChannelHyperledgerService
}

func MakeChannelController(
	channelService *services.ChannelHyperledgerService,
) *ChannelController {
	return &ChannelController{
		channelService: channelService,
	}
}

func (c *ChannelController) GetChannels() ([]string, error) {
	resultChan := make(chan services.GetChannelsResult)
	go c.channelService.GetChannels(resultChan)
	result := <-resultChan
	return result.Channels, result.Error.Error
}
