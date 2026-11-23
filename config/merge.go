package config

import "fmt"

func defaultMerge(dst, src any) error {
	dstMap, ok := dst.(*map[string]any)
	if !ok {
		return fmt.Errorf("config: merge dst must be *map[string]interface{}, got %T", dst)
	}
	srcMap, ok := convertMap(src).(map[string]any)
	if !ok {
		return fmt.Errorf("config: merge src must be map[string]interface{}, got %T", src)
	}
	if *dstMap == nil {
		*dstMap = make(map[string]any, len(srcMap))
	}
	mergeMap(*dstMap, srcMap)
	return nil
}

func mergeMap(dst, src map[string]any) {
	for key, srcValue := range src {
		if srcMap, ok := srcValue.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				mergeMap(dstMap, srcMap)
				continue
			}
		}
		dst[key] = cloneMergeValue(srcValue)
	}
}

func cloneMergeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		cloned := make(map[string]any, len(val))
		mergeMap(cloned, val)
		return cloned
	case []any:
		cloned := make([]any, len(val))
		for i, item := range val {
			cloned[i] = cloneMergeValue(item)
		}
		return cloned
	default:
		return val
	}
}
