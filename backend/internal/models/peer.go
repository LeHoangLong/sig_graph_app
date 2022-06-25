package models

type PeerId int
type Peer struct {
	Id        PeerId
	Alias     string
	PublicKey []PublicKey
}

func MakePeer(
	iId PeerId,
	iAlias string,
	iPublicKey []PublicKey,
) Peer {
	return Peer{
		Id:        iId,
		Alias:     iAlias,
		PublicKey: iPublicKey,
	}
}
