package models

type Peer struct {
	Id        int
	Alias     string
	PublicKey []PublicKey
}

func MakePeer(
	iId int,
	iAlias string,
	iPublicKey []PublicKey,
) Peer {
	return Peer{
		Id:        iId,
		Alias:     iAlias,
		PublicKey: iPublicKey,
	}
}
