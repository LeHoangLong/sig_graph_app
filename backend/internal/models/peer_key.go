package models

type PeerKey struct {
	PublicKey
}

func MakePeerKey(
	iPublicKey PublicKey,
) PeerKey {
	return PeerKey{
		PublicKey: iPublicKey,
	}
}
