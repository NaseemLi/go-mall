package set

type MYMAP[K comparable, V any] map[K]V

// DiffArray 求两个切片的差集
func DiffArray[T comparable](a []T, b []T) []T {
	var diffArray []T
	var temp MYMAP[T, struct{}] = map[T]struct{}{}
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}

// IntersectArray 求两个切片的交集
func IntersectArray[T comparable](a []T, b []T) []T {
	var inter []T
	var mp MYMAP[T, bool] = map[T]bool{}

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}

// RemoveRepByMap 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap[T comparable](slc []T) []T {
	var result []T
	var tempMap MYMAP[T, bool] = map[T]bool{}
	for _, e := range slc {
		if _, ok := tempMap[e]; !ok {
			result = append(result, e)
			tempMap[e] = true
		}
	}
	return result
}
