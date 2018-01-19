package slices

import "sort"

// Contains returns true if `needle` is found inside slice `haystack`
func Contains(haystack []interface{}, needle interface{}) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// ContainsString is a string-optimized version of Contains
func ContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// FindDuplicateStrings is a string-optimized version of FindDuplicateValues
func FindDuplicateStrings(s []string) map[string]uint {
	var duplicates = map[string]uint{}
	var previous = ""
	var count = uint(0)
	sort.Strings(s)
	for _, v := range s {
		if v != previous {
			previous = v
			count = 1
		} else {
			count++
			if count > 1 {
				duplicates[v] = count
			}
		}
	}
	return duplicates
}

// FindDuplicateValues returns a slice of values that are duplicated in the slice
func FindDuplicateValues(i []interface{}) map[interface{}]uint {
	var allValues = map[interface{}]uint{}
	for _, v := range i {
		allValues[v]++
	}
	var duplicateValues = make(map[interface{}]uint)
	for v, c := range allValues {
		if c > 1 {
			duplicateValues[v] = c
		}
	}
	return duplicateValues
}

// ListDistinctStrings is a string-optimized version of ListDistinctValues
func ListDistinctStrings(s []string) []string {
	var uniqueValues = []string{}
	var previous = ""
	sort.Strings(s)
	if len(s) > 0 && s[0] == "" {
		uniqueValues = append(uniqueValues, "")
	}
	for _, v := range s {
		if v != previous {
			uniqueValues = append(uniqueValues, v)
			previous = v
		}
	}
	return uniqueValues
}

// ListDistinctValues lists each values within a slice once and only once
func ListDistinctValues(i []interface{}) []interface{} {
	var returnSlice = []interface{}{}
	var uniqueMap = make(map[interface{}]bool)
	for _, v := range i {
		if !uniqueMap[v] {
			returnSlice = append(returnSlice, v)
			uniqueMap[v] = true
		}
	}
	return returnSlice
}
