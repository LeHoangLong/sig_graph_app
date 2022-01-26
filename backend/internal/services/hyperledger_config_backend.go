package services

type HyperledgerConfigBackend struct {
	data map[string]interface{}
}

func MakeHyperledgerConfigBackend(
	data map[string]interface{},
) *HyperledgerConfigBackend {
	return &HyperledgerConfigBackend{
		data: data,
	}
}

func (s *HyperledgerConfigBackend) Lookup(key string) (interface{}, bool) {
	if val, ok := s.data["key"]; ok {
		return val, true
	} else {
		return nil, false
	}
}
