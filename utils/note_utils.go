package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ToFixedCustom(comma string, v float64) float64 {
	// FLOAT TO STRING FORMATING DIGIT DECIMAL
	res := fmt.Sprintf("%."+comma+"f", v)
	// STRING TO FLOAT64
	v, _ = strconv.ParseFloat(res, 64)

	return v
}

func ToFixedString(comma string, v float64) string {
	// FLOAT TO STRING FORMATING DIGIT DECIMAL
	res := fmt.Sprintf("%."+comma+"f", v)

	return res
}

func FloatToAmountStringCustom(f float64, comma string) string {

	fString := ToFixedString(comma, f)
	splitDecimals := strings.Split(fString, ".")
	tempCount := len(splitDecimals[0]) % 3

	if tempCount == 0 {
		tempCount = 3
	}

	ret := make([]byte, 0)

	for i := 0; i < len(splitDecimals[0]); i++ {
		ret = append(ret, splitDecimals[0][i])
		tempCount = tempCount - 1
		if tempCount == 0 && i != len(splitDecimals[0])-1 {
			ret = append(ret, ',')
			tempCount = 3
		}
	}

	splitDecimals[0] = string(ret)

	return strings.Join(splitDecimals, ".")
}

func RenderRow(row map[string][][]byte) map[string][]string {
	ret := make(map[string][]string)

	max := 0
	for _, v := range row {
		if len(v) > max {
			max = len(v)
		}
	}

	for k, v := range row {
		ret[k] = make([]string, max)
		for i := 0; i < max; i++ {
			if len(v)-i > 0 {
				ret[k][i] = string(v[i])
			} else {
				ret[k][i] = ""
			}
		}
	}
	return ret
}
