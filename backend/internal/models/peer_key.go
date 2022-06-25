package models

type PeerKey struct {
	PublicKey
	UserId UserId
	PeerId *int
}

func MakePeerKey(
	iPublicKey PublicKey,
	iUserId UserId,
	iPeerId *int,
) PeerKey {
	return PeerKey{
		PublicKey: iPublicKey,
		UserId:    iUserId,
		PeerId:    iPeerId,
	}
}
