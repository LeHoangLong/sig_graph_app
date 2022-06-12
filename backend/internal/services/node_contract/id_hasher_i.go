package node_contract

type IdHasherI interface {
	Hash(iId string) string
}
