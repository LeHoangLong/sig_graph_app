package peer_material_services

import (
	"backend/internal/common"
	"backend/internal/models"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerMaterialClientServiceFactory struct {
	mtx             *sync.Mutex
	grpcConnections map[string]*grpc.ClientConn
}

func MakePeerMaterialClientServiceFactory() PeerMaterialClientServiceFactory {
	return PeerMaterialClientServiceFactory{
		mtx:             &sync.Mutex{},
		grpcConnections: map[string]*grpc.ClientConn{},
	}
}

func (f *PeerMaterialClientServiceFactory) BuildPeerMaterialClientService(
	iEndpoint models.PeerEndpoint,
) (PeerMaterialClientServiceI, error) {
	if iEndpoint.Protocol == models.Grpc {
		f.mtx.Lock()
		defer f.mtx.Unlock()
		if _, ok := f.grpcConnections[iEndpoint.Url]; !ok {
			conn, err := grpc.Dial(iEndpoint.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			f.grpcConnections[iEndpoint.Url] = conn
		}
		service := MakePeerMaterialClientServiceGrpc(f.grpcConnections[iEndpoint.Url])
		return &service, nil
	} else {
		return nil, common.Unsupported
	}
}
