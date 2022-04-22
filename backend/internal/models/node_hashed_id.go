package models

type NodeHashedId struct {
	Id       int
	HashedId string `json:"HashedId"`
	RawId    string `json:"RawId"`
}

func MakeNodeHashedId(
	iId int,
	iHashedId string,
	iRawId string,
) NodeHashedId {
	return NodeHashedId{
		Id:       iId,
		HashedId: iHashedId,
		RawId:    iRawId,
	}
}
