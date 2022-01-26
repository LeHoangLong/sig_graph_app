package services

import (
	"backend/internal/common"
	"backend/internal/repositories"
)

type ChannelStorageService struct {
	repository repositories.ChannelRepositoryI
}

func (s *ChannelStorageService) AddChannel(iChannelName string, oResult chan common.WrappedError) {
	go s.repository.AddChannel(iChannelName, oResult)
}
