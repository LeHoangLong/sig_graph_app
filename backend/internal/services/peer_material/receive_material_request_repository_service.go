package peer_material_services

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/repositories"
	peer_material_repositories "backend/internal/repositories/peer_material"
	"context"
)

type ReceiveMaterialRequestRepositoryService struct {
	materialRepository               repositories.MaterialRepositoryI
	receiveMaterialRequestRepository peer_material_repositories.MaterialReceiveRequestRepositoryI
	peerProtocolRepository           repositories.PeerProtocolRepositoryI
}

func MakeReceiveMaterialRequestRepositoryService(
	iMaterialRepository repositories.MaterialRepositoryI,
	iReceiveMaterialRequestRepository peer_material_repositories.MaterialReceiveRequestRepositoryI,
	iPeerProtocolRepository repositories.PeerProtocolRepositoryI,
) *ReceiveMaterialRequestRepositoryService {
	return &ReceiveMaterialRequestRepositoryService{
		materialRepository:               iMaterialRepository,
		receiveMaterialRequestRepository: iReceiveMaterialRequestRepository,
		peerProtocolRepository:           iPeerProtocolRepository,
	}
}

func (s *ReceiveMaterialRequestRepositoryService) CreateOutboundReceiveMaterialRequest(
	iContext context.Context,
	iRecipientPeerId models.PeerId,
	iSenderUserId models.UserId,
	iSenderPublicKeyId models.PublicKeyId,
	iMainMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iTransferTime models.CustomTime,
) (models.OutboundMaterialReceiveRequest, error) {
	return s.receiveMaterialRequestRepository.CreateOutboundReceiveMaterialRequest(
		iContext,
		iRecipientPeerId,
		iSenderUserId,
		iSenderPublicKeyId,
		iMainMaterial,
		iRelatedMaterials,
		iOptions,
		iTransferTime,
		models.PENDING,
	)
}

func (s *ReceiveMaterialRequestRepositoryService) CreateInboundReceiveMaterialRequest(
	iContext context.Context,
	iRecipientUserId models.UserId,
	iSenderPublicKeyId models.PublicKeyId,
	iMainMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iTransferTime models.CustomTime,
	iSenderEndpoints []SenderEndpoint,
) (models.InboundMaterialReceiveRequest, error) {
	if len(iSenderEndpoints) == 0 {
		return models.InboundMaterialReceiveRequest{}, common.InvalidArgument
	}
	protocolWithMajorVersions := map[repositories.ProtocolWithMajorVersion]bool{}
	for i := range iSenderEndpoints {
		protocolWithMajorVersion := repositories.MakeProtocolWithMajorVersion(
			iSenderEndpoints[i].Protocol,
			iSenderEndpoints[i].MajorVersion,
		)
		protocolWithMajorVersions[protocolWithMajorVersion] = true
	}
	supportedProtocols, err := s.peerProtocolRepository.FilterSupportedProtocolsWithMajorVersion(
		iContext,
		protocolWithMajorVersions,
	)

	if err != nil {
		return models.InboundMaterialReceiveRequest{}, err
	}

	senderEndpoints := []models.SenderEndpoint{}
	for i := range iSenderEndpoints {
		protocolWithMajorVersion := repositories.MakeProtocolWithMajorVersion(
			iSenderEndpoints[i].Protocol,
			iSenderEndpoints[i].MajorVersion,
		)
		if protocol, ok := supportedProtocols[protocolWithMajorVersion]; ok {
			/// find the lowest minor version that is at least equals to the endpoint minor
			lowestMinorIndex := -1
			lowestMinor := -1
			for i := range protocol {
				if protocol[i].MinorVersion >= iSenderEndpoints[i].MinorVersion {
					if protocol[i].MinorVersion < lowestMinor || lowestMinorIndex == -1 {
						lowestMinor = protocol[i].MinorVersion
						lowestMinorIndex = i
					}
				}
			}

			if lowestMinorIndex >= 0 {
				senderEndpoints = append(senderEndpoints, models.MakeSenderEndpoint(
					protocol[lowestMinorIndex],
					iSenderEndpoints[i].Url,
				))
			}
		}
	}

	if len(senderEndpoints) == 0 {
		return models.InboundMaterialReceiveRequest{}, common.Unsupported
	}

	return s.receiveMaterialRequestRepository.CreateInboundReceiveMaterialRequest(
		iContext,
		iRecipientUserId,
		iSenderPublicKeyId,
		iMainMaterial,
		iRelatedMaterials,
		iOptions,
		iTransferTime,
		models.PENDING,
		senderEndpoints,
	)
}

func (s *ReceiveMaterialRequestRepositoryService) UpdateMaterialReceiveRequestStatus(
	iContext context.Context,
	iRequestId models.MaterialReceiveRequestId,
	iStatus models.MaterialReceiveRequestStatus,
) error {
	err := s.receiveMaterialRequestRepository.UpdateMaterialReceiveRequestStatus(
		iContext,
		iRequestId,
		iStatus,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *ReceiveMaterialRequestRepositoryService) FetchInboundReceiveMaterialRequestsByUser(
	iContext context.Context,
	iUserId models.UserId,
	iStatus []models.MaterialReceiveRequestStatus,
) ([]models.InboundMaterialReceiveRequest, error) {
	simplifiedRequests, err := s.receiveMaterialRequestRepository.FetchInboundReceiveMaterialRequestsByUserId(
		iContext,
		iUserId,
		iStatus,
	)

	if err != nil {
		return []models.InboundMaterialReceiveRequest{}, err
	}

	requests, err := s.makeMaterialReceiveRequestsFromSimplifiedRequests(
		iContext,
		simplifiedRequests,
	)

	inboundRequests := []models.InboundMaterialReceiveRequest{}
	for i := range requests {
		inboundRequest := models.MakeInboundMaterialReceiveRequest(
			requests[i],
			simplifiedRequests[i].RecipientUserId,
			simplifiedRequests[i].SenderPublicKey,
			simplifiedRequests[i].SenderEndpoints,
		)
		inboundRequests = append(inboundRequests, inboundRequest)
	}

	if err != nil {
		return []models.InboundMaterialReceiveRequest{}, err
	}

	return inboundRequests, nil
}

func (s *ReceiveMaterialRequestRepositoryService) makeMaterialReceiveRequestsFromSimplifiedRequests(
	iContext context.Context,
	simplifiedRequests []peer_material_repositories.SimplifiedInboundReceiveMaterialRequest,
) ([]models.MaterialReceiveRequest, error) {
	materialsIdToFetch := map[models.NodeId]bool{}
	for _, request := range simplifiedRequests {
		materialsIdToFetch[request.ToBeReceivedMaterialId] = true
		for _, relatedMaterialId := range request.RelatedMaterialIds {
			materialsIdToFetch[relatedMaterialId] = true
		}
	}

	materials, err := s.materialRepository.FetchMaterialsById(
		iContext,
		materialsIdToFetch,
	)

	if err != nil {
		return []models.MaterialReceiveRequest{}, err
	}

	requests := []models.MaterialReceiveRequest{}
	for _, request := range simplifiedRequests {
		relatedMaterials := []models.Material{}
		for _, id := range request.RelatedMaterialIds {
			relatedMaterials = append(relatedMaterials, materials[id])
		}

		receiveMaterialRequest := models.MakeMaterialReceiveRequest(
			request.Id,
			materials[request.ToBeReceivedMaterialId],
			relatedMaterials,
			request.SignatureOptions,
			request.TransferTime,
			request.Status,
		)

		requests = append(requests, receiveMaterialRequest)
	}
	return requests, nil
}
