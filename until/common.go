package until

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetStructTagContent(i interface{}, fieldNname string, tagName string) (string, bool) {
	typeOfI := reflect.TypeOf(i)
	for typeOfI.Kind() == reflect.Ptr {
		typeOfI = typeOfI.Elem()
	}
	if iType, ok := typeOfI.FieldByName(fieldNname); ok {
		return iType.Tag.Lookup(tagName)
	}
	return "", false
}

// 获取正在运行的函数信息
func GetFuncInfo(skip int) (funcName string, fileName string, line int, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	return f.Name(), file, line, ok
}

// 格式化时间为日期
// 假设当前日期：2020-08-24
//
//		2020-08-24 18:20 => 18:20
//	 2020-08-23 18:20 => 昨天18:20
//	 2020-08-22 18:20 => 前天18:20
//	 2020-08-21 18:20 => 08-21
//	 2019-08-21 18:20 => 2019-08-21
func formatTime(t time.Time, baseUnix int64) (date string) {
	base := time.Unix(baseUnix, 0)
	cur := time.Date(base.Year(), base.Month(), base.Day()+1, 0, 0, 0, 0, time.Local)
	diff := (int)((cur.Unix() - t.Unix()) / 86400)
	if diff < 0 {
		diff = 0
	}
	format := "2006-01-02" // YYYY-MM-DD
	diffForFormat := [3]string{
		"15:04",   // 当天：HH:mm
		"昨天15:04", // 前一天：昨天HH:mm
		"前天15:04", // 前两天：前天HH:mm
	}
	if diff < len(diffForFormat) {
		format = diffForFormat[diff]
	} else {
		if cur.Year() == t.Year() { // 当年：MM-DD
			format = "01-02"
		}
	}
	date = t.Format(format)
	return
}

func FormatTime(unixTime int64) (date string) {
	date = formatTime(time.Unix(unixTime, 0), time.Now().Unix())
	return
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// map数据转成struct结构
func Map2Struct(mapData interface{}, structData interface{}) error {
	retBuffer, err := json.Marshal(mapData)
	if err == nil {
		err = json.Unmarshal(retBuffer, structData)
	}
	return nil
}

// 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int64) []int64 {
	var result []int64
	tempMap := map[int64]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func StringInSlice(s string, list []string) bool {
	for _, val := range list {
		if val == s {
			return true
		}
	}
	return false
}

// 取交集
func Intersect4Int64(big, small []int64) (out []int64) {
	amap := make(map[int64]int, len(big))
	for _, v := range big {
		amap[v] = 1
	}
	for _, v := range small {
		if _, ok := amap[v]; ok {
			out = append(out, v)
		}
	}
	return
}

// 取int64切片的差集
func NonIntersect4Int64(o, n []int64) []int64 {
	if len(o) <= 0 {
		return n
	}
	if len(n) <= 0 {
		return o
	}
	m := make(map[int64]int)
	var arr []int64
	for _, v := range o {
		m[v]++
	}
	for _, v := range n {
		m[v]++
	}
	for k, v := range m {
		if v == 1 {
			arr = append(arr, k)
		}
	}
	return arr
}

// 取字符切片的差集
func NonIntersect4String(o, n []string) []string {
	if len(o) <= 0 {
		return n
	}
	if len(n) <= 0 {
		return o
	}
	m := make(map[string]int)
	var arr []string
	for _, v := range o {
		m[v]++
	}
	for _, v := range n {
		m[v]++
	}
	for k, v := range m {
		if v == 1 {
			arr = append(arr, k)
		}
	}
	return arr
}

// 将字符串分割成int64的切片
func StringSplitInt64Slice(str, sep string) []int64 {
	ids := []int64{}
	if str == "" {
		return ids
	}

	strList := strings.Split(strings.Trim(str, sep), sep)
	if len(strList) == 0 {
		return ids
	}
	for _, val := range strList {
		if val == "" {
			continue
		}
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, int64(id))
	}
	return ids
}

// map/struc转json字符串
func Params2JsonStr(params interface{}) string {
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	return string(paramsStr)
}

func Json2Map(params string) map[string]interface{} {
	result := make(map[string]interface{})
	_ = json.Unmarshal([]byte(params), &result)
	return result
}

func StringDate2Time(layout string, t string) (*time.Time, error) {
	temp, err := time.Parse(layout, t)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

// 生成更新参数
func MakeUpdateParams(Type string) map[string]interface{} {
	updateFields := make(map[string]interface{})
	if Type == "DELETE" {
		updateFields["is_delete"] = 1
		updateFields["delete_time"] = time.Now().Unix()
	} else {
		updateFields["update_time"] = time.Now().Unix()
	}
	return updateFields
}

// 获取当前时间time类型
func LocalTime() time.Time {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	lastTime := strings.ReplaceAll(nowTime[0:19], "T", " ")
	localTime, _ := time.ParseInLocation("2006-01-02 15:04:05", lastTime, time.Local)
	return localTime
}

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	// 重写time转换成json之后的格式
	stmp := fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02 15:04:05"))
	fmt.Println(stmp)
	return []byte(stmp), nil
}
