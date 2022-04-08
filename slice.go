package golibrary

// 交集 Set
func IntIntersect(data1 []int, data2 []int) (data3 []int) {
	if len(data1) == 0 {
		return data2
	}
	if len(data2) == 0 {
		return data1
	}
	dict := map[int]int{}
	for _, v := range data1 {
		if _, ok := dict[v]; ok {
			dict[v] += 1
		} else {
			dict[v] = 1
		}
	}
	for _, v := range data2 {
		if _, ok := dict[v]; ok {
			dict[v] += 1
		} else {
			dict[v] = 1
		}
	}
	for i, v := range dict {
		if v > 1 {
			data3 = append(data3, i)
		}
	}
	return
}

// 移除数组中重复的元素
func Int64Unique(data1 []int64) (rst []int64) {
	data := make(map[int64]bool, len(data1))
	for _, v := range data1 {
		if _, ok := data[v]; !ok {
			data[v] = true
		}
	}
	for i, v := range data {
		if v {
			rst = append(rst, i)
		}
	}
	return
}

func IntUnique(data1 []int) (rst []int) {
	data := make(map[int]bool, len(data1))
	for _, v := range data1 {
		if _, ok := data[v]; !ok {
			data[v] = true
		}
	}
	for i, v := range data {
		if v {
			rst = append(rst, i)
		}
	}
	return
}

// 返回索引
func IndexOfInt64(slice []int64, find int64) int {
	for i, v := range slice {
		if v == find {
			return i
		}
	}
	return -1
}


func IndexOfInt(slice []int, find int) int {
	for i, v := range slice {
		if v == find {
			return i
		}
	}
	return -1
}

// 求两个切片的差集, data1 里不在 data2 中的元素
func DiffInt64(data1 []int64, data2 []int64) (res []int64) {
	for i := range data1 {
		if IndexOfInt64(data2, data1[i]) < 0 {
			res = append(res, data1[i])
		}
	}
	return
}
