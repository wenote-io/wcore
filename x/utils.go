package x

// RemoveRepByMap int 去除重复
func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]struct{}{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = struct{}{}
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
