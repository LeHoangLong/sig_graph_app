package peer_material_services

import "backend/internal/models"

func GenerateNodesFromMaterials(
	iMaterials []models.Material,
) []Node {
	temp := map[int]models.Material{}
	for _, material := range iMaterials {
		temp[*material.Id] = material
	}

	ret := make([]Node, 0, len(iMaterials))
	for _, material := range iMaterials {
		children := []string{}
		for id := range material.ChildrenIds {
			if material, ok := temp[id]; ok {
				children = append(children, material.NodeId)
			}
		}
		parents := []string{}
		for id := range material.ParentIds {
			if material, ok := temp[id]; ok {
				parents = append(parents, material.NodeId)
			}
		}
		ret = append(ret, MakeNode(children, parents))
	}

	return ret
}
