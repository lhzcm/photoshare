package utility

import "testing"

func TestStringToIntArray(t *testing.T) {
	str := "32,34,23,12,343,213,34"
	sucarry := []int{32, 34, 23, 12, 343, 213, 34}
	array, err := StringToIntArray(str, ",")
	t.Log(array)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if len(sucarry) != len(array) {
		t.Log("长度有误")
		t.Fail()
	}
	for i, item := range sucarry {
		if item != array[i] {
			t.Log("结果有误")
			t.Fail()
		}
	}
	if array, err = StringToIntArray("12,s,er,sa,", ","); err == nil {
		t.Log("错误结果有误")
		t.Fail()
	}
}
