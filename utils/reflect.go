package utils

import (
	"fmt"
	"reflect"
)

func MapChanges(map1, map2 map[string]interface{}, prefix string, modified map[string][2]interface{}) {
	for key, val1 := range map1 {
		val2, exists := map2[key]
		if !exists {
			continue // Ignore removed keys
		}

		// If both values are maps, recurse
		subMap1, ok1 := val1.(map[string]interface{})
		subMap2, ok2 := val2.(map[string]interface{})

		if ok1 && ok2 {
			MapChanges(subMap1, subMap2, key+".", modified)
			continue
		}

		// Handle Array value
		subMapArr1, ok1 := val1.([]interface{})
		subMapArr2, ok2 := val2.([]interface{})

		if ok1 && ok2 {
			for i := 0; i < min(len(subMapArr1), len(subMapArr2)); i++ {
				map1Item, ok1 := subMapArr1[i].(map[string]interface{})
				map2Item, ok2 := subMapArr2[i].(map[string]interface{})
				if ok1 && ok2 {
					MapChanges(map1Item, map2Item, fmt.Sprintf("%s%d.", key+".", i+1), modified)
				} else if !reflect.DeepEqual(subMapArr1[i], subMapArr2[i]) {
					modified[fmt.Sprintf("%s%d", key, i)] = [2]interface{}{subMapArr1[i], subMapArr2[i]}
				}
			}

			if len(subMapArr2) > len(subMapArr1) {
				for i := len(subMapArr1); i < len(subMapArr2); i++ {
					modified[fmt.Sprintf("%s%d", key+".", i+1)] = [2]interface{}{nil, subMapArr2[i]}
				}
			}

			continue
		}

		// If values are different, add to modified result
		if !reflect.DeepEqual(val1, val2) {
			modified[prefix+key] = [2]interface{}{val1, val2}
		}
	}
}

// Helper function to get the minimum of two values
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
