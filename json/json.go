package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// JSONUtil 提供了一系列 JSON 相关的工具方法
type JSONUtil struct{}

// Marshal 将对象转换为 JSON 字节切片
func (ju *JSONUtil) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent 将对象转换为格式化的 JSON 字节切片
func (ju *JSONUtil) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal 将 JSON 字节切片解析为对象
func (ju *JSONUtil) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ToJSON 将对象转换为 JSON 字符串
func (ju *JSONUtil) ToJSON(v interface{}) (string, error) {
	b, err := ju.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON 将 JSON 字符串解析为对象
func (ju *JSONUtil) FromJSON(jsonStr string, v interface{}) error {
	return ju.Unmarshal([]byte(jsonStr), v)
}

// ToJSONFile 将对象保存为 JSON 文件
func (ju *JSONUtil) ToJSONFile(v interface{}, filename string) error {
	data, err := ju.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// FromJSONFile 从 JSON 文件读取并解析为对象
func (ju *JSONUtil) FromJSONFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return ju.Unmarshal(data, v)
}

// IsValidJSON 检查字符串是否为有效的 JSON
func (ju *JSONUtil) IsValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// PrettyPrint 返回格式化的 JSON 字符串
func (ju *JSONUtil) PrettyPrint(v interface{}) (string, error) {
	b, err := ju.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compact 返回压缩的 JSON 字符串
func (ju *JSONUtil) Compact(v interface{}) (string, error) {
	b, err := ju.Marshal(v)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, b); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// DeepCopy 深拷贝 JSON 对象
func (ju *JSONUtil) DeepCopy(src, dst interface{}) error {
	data, err := ju.Marshal(src)
	if err != nil {
		return err
	}
	return ju.Unmarshal(data, dst)
}

// MergeJSON 合并两个 JSON 对象
func (ju *JSONUtil) MergeJSON(json1, json2 string) (string, error) {
	var m1, m2, result map[string]interface{}

	if err := ju.FromJSON(json1, &m1); err != nil {
		return "", err
	}
	if err := ju.FromJSON(json2, &m2); err != nil {
		return "", err
	}

	result = ju.mergeMap(m1, m2)
	return ju.ToJSON(result)
}

// mergeMap 递归合并两个 map
func (ju *JSONUtil) mergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		if v2, ok := result[k]; ok {
			if vMap, ok := v.(map[string]interface{}); ok {
				if v2Map, ok := v2.(map[string]interface{}); ok {
					result[k] = ju.mergeMap(v2Map, vMap)
					continue
				}
			}
		}
		result[k] = v
	}

	return result
}

// GetValueByPath 通过路径获取 JSON 中的值
func (ju *JSONUtil) GetValueByPath(jsonStr, path string) (interface{}, error) {
	var data interface{}
	if err := ju.FromJSON(jsonStr, &data); err != nil {
		return nil, err
	}

	keys := strings.Split(path, ".")
	return ju.getValueByKeys(data, keys)
}

// getValueByKeys 递归获取嵌套值
func (ju *JSONUtil) getValueByKeys(data interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return data, nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		if value, ok := v[keys[0]]; ok {
			return ju.getValueByKeys(value, keys[1:])
		}
		return nil, fmt.Errorf("key not found: %s", keys[0])
	case []interface{}:
		if keys[0] == "*" {
			var results []interface{}
			for _, item := range v {
				result, err := ju.getValueByKeys(item, keys[1:])
				if err == nil {
					results = append(results, result)
				}
			}
			return results, nil
		}
		return nil, fmt.Errorf("invalid key for array: %s", keys[0])
	default:
		return nil, errors.New("invalid JSON structure")
	}
}

// Diff 比较两个 JSON 对象，返回差异
func (ju *JSONUtil) Diff(json1, json2 string) (map[string]interface{}, error) {
	var obj1, obj2 interface{}
	if err := ju.FromJSON(json1, &obj1); err != nil {
		return nil, err
	}
	if err := ju.FromJSON(json2, &obj2); err != nil {
		return nil, err
	}

	diff := make(map[string]interface{})
	ju.diffObjects("", obj1, obj2, diff)
	return diff, nil
}

// diffObjects 递归比较两个对象
func (ju *JSONUtil) diffObjects(prefix string, obj1, obj2 interface{}, diff map[string]interface{}) {
	switch v1 := obj1.(type) {
	case map[string]interface{}:
		v2, ok := obj2.(map[string]interface{})
		if !ok {
			diff[prefix] = map[string]interface{}{"old": obj1, "new": obj2}
			return
		}
		for k, v := range v1 {
			newPrefix := prefix
			if newPrefix != "" {
				newPrefix += "."
			}
			newPrefix += k
			if v2v, ok := v2[k]; ok {
				ju.diffObjects(newPrefix, v, v2v, diff)
			} else {
				diff[newPrefix] = map[string]interface{}{"old": v, "new": nil}
			}
		}
		for k, v := range v2 {
			if _, ok := v1[k]; !ok {
				newPrefix := prefix
				if newPrefix != "" {
					newPrefix += "."
				}
				newPrefix += k
				diff[newPrefix] = map[string]interface{}{"old": nil, "new": v}
			}
		}
	case []interface{}:
		v2, ok := obj2.([]interface{})
		if !ok {
			diff[prefix] = map[string]interface{}{"old": obj1, "new": obj2}
			return
		}
		if len(v1) != len(v2) {
			diff[prefix] = map[string]interface{}{"old": obj1, "new": obj2}
			return
		}
		for i := range v1 {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			ju.diffObjects(newPrefix, v1[i], v2[i], diff)
		}
	default:
		if !reflect.DeepEqual(obj1, obj2) {
			diff[prefix] = map[string]interface{}{"old": obj1, "new": obj2}
		}
	}
}

// StreamingDecode 流式解码大型 JSON 文件
func (ju *JSONUtil) StreamingDecode(r io.Reader, callback func(json.Token) error) error {
	dec := json.NewDecoder(r)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := callback(t); err != nil {
			return err
		}
	}
	return nil
}

// New 创建一个新的 JSONUtil 实例
func New() *JSONUtil {
	return &JSONUtil{}
}
