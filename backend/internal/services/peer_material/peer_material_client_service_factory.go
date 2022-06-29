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
	iProtocolName models.ProtocolName,
	iMajorVersion int,
	iMinorVersion int,
	iUrl string,
) (PeerMaterialClientServiceI, error) {
	if iProtocolName == models.GRPC {
		f.mtx.Lock()
		defer f.mtx.Unlock()
		if _, ok := f.grpcConnections[iUrl]; !ok {
			conn, err := grpc.Dial(iUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			f.grpcConnections[iUrl] = conn
		}
		service := MakePeerMaterialClientServiceGrpc(f.grpcConnections[iUrl])
		return &service, nil
	} else {
		return nil, common.Unsupported
	}
}
