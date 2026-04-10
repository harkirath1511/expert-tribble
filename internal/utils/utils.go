package utils

func GetMap(m map[string]any, key string) (map[string]any, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}
	child, ok := v.(map[string]any)
	return child, ok
}

func GetFloat(m map[string]any, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	n, ok := v.(float64)
	return n, ok
}
