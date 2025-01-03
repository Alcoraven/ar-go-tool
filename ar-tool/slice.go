package artool

import "strconv"

func InSlice[T int | int32 | int64 | uint | uint32 | uint64 | string](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func SliceIntToString[T int | int32 | int64](arr []T) []string {
	var res []string
	for _, v := range arr {
		res = append(res, strconv.FormatInt(int64(v), 10))
	}
	return res
}

func SliceUintToString[T uint | uint32 | uint64](arr []T) []string {
	var res []string
	for _, v := range arr {
		res = append(res, strconv.FormatUint(uint64(v), 10))
	}
	return res
}

func SliceStringToUint(arr []string) []uint64 {
	var res []uint64
	for _, v := range arr {
		d, _ := strconv.ParseUint(v, 10, 64)
		res = append(res, d)
	}
	return res
}

func SliceStringToInt(arr []string) []int64 {
	var res []int64
	for _, v := range arr {
		d, _ := strconv.ParseInt(v, 10, 64)
		res = append(res, d)
	}
	return res
}

func SliceReverse[T int | int32 | int64 | uint | uint32 | uint64 | string](list *[]T) {
	length := len(*list)
	for i := 0; i < length/2; i += 1 {
		(*list)[i], (*list)[length-1-i] = (*list)[length-1-i], (*list)[i]
	}
}
