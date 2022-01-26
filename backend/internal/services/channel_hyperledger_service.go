package services

import (
	"backend/internal/common"
	"os/exec"
	"strings"
)

type ChannelHyperledgerService struct {
}

type GetChannelsResult struct {
	Channels []string
	Error    common.WrappedError
}

func MakeChannelHyperledgerService() *ChannelHyperledgerService {
	return &ChannelHyperledgerService{}
}

func (s *ChannelHyperledgerService) GetChannels(oResult chan GetChannelsResult) {
	cmd := exec.Command("./tools/bin/peer", "channel list")
	stdout, err := cmd.Output()

	if err != nil {
		oResult <- GetChannelsResult{
			Error: common.WrappedError{
				Code:  common.Unknown,
				Error: err,
			},
		}
		return
	}

	oResult <- GetChannelsResult{
		Channels: strings.Split(string(stdout), "\n"),
	}

}
