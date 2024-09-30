package randutil

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"
)

// 定义常量
const (
	// 默认字符集
	defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 默认数字字符集
	defaultDigits = "0123456789"
	// 默认字母字符集
	defaultLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Random 结构体用于生成随机数和字符串
type Random struct {
	charset string
}

// New 创建一个新的 Random 实例
func New() *Random {
	return &Random{charset: defaultCharset}
}

// SetCharset 设置自定义字符集
func (r *Random) SetCharset(charset string) {
	r.charset = charset
}

// Int 生成指定范围内的随机整数 [min, max]
func (r *Random) Int(min, max int) (int, error) {
	if min > max {
		min, max = max, min
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + min, nil
}

// Float64 生成指定范围内的随机浮点数 [min, max)
func (r *Random) Float64(min, max float64) (float64, error) {
	if min >= max {
		return 0, errors.New("min should be less than max")
	}

	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}

	// Convert to float64 in [0, 1)
	f := float64(binary.LittleEndian.Uint64(b[:])) / (1 << 64)

	return min + f*(max-min), nil
}

// String 生成指定长度的随机字符串
func (r *Random) String(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = r.charset[b%byte(len(r.charset))]
	}
	return string(bytes), nil
}

// Digits 生成指定长度的随机数字字符串
func (r *Random) Digits(length int) (string, error) {
	oldCharset := r.charset
	r.charset = defaultDigits
	defer func() { r.charset = oldCharset }()
	return r.String(length)
}

// Letters 生成指定长度的随机字母字符串
func (r *Random) Letters(length int) (string, error) {
	oldCharset := r.charset
	r.charset = defaultLetters
	defer func() { r.charset = oldCharset }()
	return r.String(length)
}

// Bytes 生成指定长度的随机字节切片
func (r *Random) Bytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Base64 生成指定长度的随机 Base64 编码字符串
func (r *Random) Base64(length int) (string, error) {
	bytes, err := r.Bytes((length*3 + 3) / 4)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded[:length], nil
}

// Hex 生成指定长度的随机十六进制字符串
func (r *Random) Hex(length int) (string, error) {
	bytes, err := r.Bytes((length + 1) / 2)
	if err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(bytes)
	return encoded[:length], nil
}

// UUID 生成 UUID v4
func (r *Random) UUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	// 设置版本 (4) 和变体 (2)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return hex.EncodeToString(uuid[:4]) + "-" +
		hex.EncodeToString(uuid[4:6]) + "-" +
		hex.EncodeToString(uuid[6:8]) + "-" +
		hex.EncodeToString(uuid[8:10]) + "-" +
		hex.EncodeToString(uuid[10:]), nil
}

// ShuffleString 随机打乱字符串
func (r *Random) ShuffleString(s string) (string, error) {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		j, err := r.Int(0, i)
		if err != nil {
			return "", err
		}
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

// Choice 从给定的切片中随机选择一个元素
func (r *Random) Choice(slice interface{}) (interface{}, error) {
	switch v := slice.(type) {
	case []string:
		if len(v) == 0 {
			return nil, nil
		}
		i, err := r.Int(0, len(v)-1)
		if err != nil {
			return nil, err
		}
		return v[i], nil
	case []int:
		if len(v) == 0 {
			return nil, nil
		}
		i, err := r.Int(0, len(v)-1)
		if err != nil {
			return nil, err
		}
		return v[i], nil
	case []float64:
		if len(v) == 0 {
			return nil, nil
		}
		i, err := r.Int(0, len(v)-1)
		if err != nil {
			return nil, err
		}
		return v[i], nil
	default:
		return nil, nil
	}
}

// WeightedChoice 根据权重从给定的切片中随机选择一个元素
func (r *Random) WeightedChoice(choices []string, weights []int) (string, error) {
	if len(choices) != len(weights) {
		return "", nil
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	n, err := r.Int(1, totalWeight)
	if err != nil {
		return "", err
	}

	for i, w := range weights {
		n -= w
		if n <= 0 {
			return choices[i], nil
		}
	}

	return "", nil
}

// Sample 从给定的切片中随机选择n个不重复的元素
func (r *Random) Sample(slice interface{}, n int) (interface{}, error) {
	switch v := slice.(type) {
	case []string:
		return r.sampleStrings(v, n)
	case []int:
		return r.sampleInts(v, n)
	case []float64:
		return r.sampleFloat64s(v, n)
	default:
		return nil, nil
	}
}

func (r *Random) sampleStrings(slice []string, n int) ([]string, error) {
	if n > len(slice) {
		n = len(slice)
	}
	result := make([]string, n)
	taken := make(map[int]bool)
	for i := 0; i < n; i++ {
		j, err := r.Int(0, len(slice)-1)
		if err != nil {
			return nil, err
		}
		for taken[j] {
			j, err = r.Int(0, len(slice)-1)
			if err != nil {
				return nil, err
			}
		}
		result[i] = slice[j]
		taken[j] = true
	}
	return result, nil
}

func (r *Random) sampleInts(slice []int, n int) ([]int, error) {
	if n > len(slice) {
		n = len(slice)
	}
	result := make([]int, n)
	taken := make(map[int]bool)
	for i := 0; i < n; i++ {
		j, err := r.Int(0, len(slice)-1)
		if err != nil {
			return nil, err
		}
		for taken[j] {
			j, err = r.Int(0, len(slice)-1)
			if err != nil {
				return nil, err
			}
		}
		result[i] = slice[j]
		taken[j] = true
	}
	return result, nil
}

func (r *Random) sampleFloat64s(slice []float64, n int) ([]float64, error) {
	if n > len(slice) {
		n = len(slice)
	}
	result := make([]float64, n)
	taken := make(map[int]bool)
	for i := 0; i < n; i++ {
		j, err := r.Int(0, len(slice)-1)
		if err != nil {
			return nil, err
		}
		for taken[j] {
			j, err = r.Int(0, len(slice)-1)
			if err != nil {
				return nil, err
			}
		}
		result[i] = slice[j]
		taken[j] = true
	}
	return result, nil
}

// RandomDate 生成指定范围内的随机日期
func (r *Random) RandomDate(start, end time.Time) (time.Time, error) {
	diff := end.Sub(start)
	randomDuration, err := r.Int64(0, int64(diff))
	if err != nil {
		return time.Time{}, err
	}
	return start.Add(time.Duration(randomDuration)), nil
}

// Int64 生成指定范围内的随机 int64 [min, max]
func (r *Random) Int64(min, max int64) (int64, error) {
	if min > max {
		min, max = max, min
	}
	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}
	return n.Int64() + min, nil
}

// Bool 生成随机布尔值
func (r *Random) Bool() (bool, error) {
	n, err := r.Int(0, 1)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// Password 生成指定长度的随机密码，包含大小写字母、数字和特殊字符
func (r *Random) Password(length int) (string, error) {
	if length < 4 {
		return "", nil
	}

	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	r.charset = lowercase + uppercase + digits + symbols

	password, err := r.String(length)
	if err != nil {
		return "", err
	}

	// 确保密码包含至少一个小写字母、一个大写字母、一个数字和一个特殊字符
	hasLower := strings.ContainsAny(password, lowercase)
	hasUpper := strings.ContainsAny(password, uppercase)
	hasDigit := strings.ContainsAny(password, digits)
	hasSymbol := strings.ContainsAny(password, symbols)

	for !(hasLower && hasUpper && hasDigit && hasSymbol) {
		password, err = r.String(length)
		if err != nil {
			return "", err
		}
		hasLower = strings.ContainsAny(password, lowercase)
		hasUpper = strings.ContainsAny(password, uppercase)
		hasDigit = strings.ContainsAny(password, digits)
		hasSymbol = strings.ContainsAny(password, symbols)
	}

	return password, nil
}

// 恢复默认字符集
func (r *Random) resetCharset() {
	r.charset = defaultCharset
}
