package models

type Node struct {
	Id              int    `json:"Id"`
	NextNodeIds     []int  `json:"NextNodeIds"`
	PreviousNodeIds []int  `json:"NextNodeIds"`
	Data            string `json:"Data"`
}
