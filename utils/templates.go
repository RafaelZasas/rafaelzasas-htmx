package utils

import "fmt"

func GetFromMap(m map[string]interface{}, key string) (any, error) {
	if value, ok := m[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("key not found")
}

func TmplSlice(elements ...any) []any {
	return elements
}

func TmplMap(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("TmplMap: odd number of arguments")
	}
	data := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("TmplMap: non-string key")
		}
		data[key] = pairs[i+1]
	}
	return data, nil
}

func WithComponentData(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("WithInputData: odd number of arguments")
	}
	data := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("WithInputData: non-string key")
		}
		data[key] = pairs[i+1]
	}
	return data, nil
}
