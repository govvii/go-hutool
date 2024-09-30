package codec

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"html"
	"io"
	"math/big"
	"net/url"
	"os"
	"strings"
)

// Base62字符集
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64URL编码
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URL解码
func Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Base32编码
func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

// Base32解码
func Base32Decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

// Base62编码
func Base62Encode(data []byte) string {
	bi := new(big.Int).SetBytes(data)
	encoded := ""
	for bi.Sign() > 0 {
		mod := new(big.Int)
		bi.DivMod(bi, big.NewInt(62), mod)
		encoded = string(base62Chars[mod.Int64()]) + encoded
	}
	return encoded
}

// Base62解码
func Base62Decode(s string) ([]byte, error) {
	bi := new(big.Int)
	for _, r := range s {
		bi.Mul(bi, big.NewInt(62))
		index := strings.IndexRune(base62Chars, r)
		bi.Add(bi, big.NewInt(int64(index)))
	}
	return bi.Bytes(), nil
}

// HexEncode 将字节切片编码为十六进制字符串
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode 将十六进制字符串解码为字节切片
func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// MD5 计算字符串的MD5哈希值
func MD5(s string) string {
	return HashString(md5.New(), s)
}

// MD5Bytes 计算字节切片的MD5哈希值
func MD5Bytes(data []byte) []byte {
	return HashBytes(md5.New(), data)
}

// SHA1 计算字符串的SHA1哈希值
func SHA1(s string) string {
	return HashString(sha1.New(), s)
}

// SHA1Bytes 计算字节切片的SHA1哈希值
func SHA1Bytes(data []byte) []byte {
	return HashBytes(sha1.New(), data)
}

// SHA256 计算字符串的SHA256哈希值
func SHA256(s string) string {
	return HashString(sha256.New(), s)
}

// SHA256Bytes 计算字节切片的SHA256哈希值
func SHA256Bytes(data []byte) []byte {
	return HashBytes(sha256.New(), data)
}

// SHA512 计算字符串的SHA512哈希值
func SHA512(s string) string {
	return HashString(sha512.New(), s)
}

// SHA512Bytes 计算字节切片的SHA512哈希值
func SHA512Bytes(data []byte) []byte {
	return HashBytes(sha512.New(), data)
}

// HashString 计算字符串的哈希值并返回十六进制字符串
func HashString(h hash.Hash, s string) string {
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// HashBytes 计算字节切片的哈希值
func HashBytes(h hash.Hash, data []byte) []byte {
	h.Write(data)
	return h.Sum(nil)
}

// HashFile 计算文件的哈希值
func HashFile(h hash.Hash, filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// ROT13 对字符串进行ROT13编码/解码
func ROT13(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			result.WriteRune('A' + (r-'A'+13)%26)
		case r >= 'a' && r <= 'z':
			result.WriteRune('a' + (r-'a'+13)%26)
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// URLEncode 对字符串进行URL编码
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode 对URL编码的字符串进行解码
func URLDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// HTMLEscape 对字符串进行HTML转义
func HTMLEscape(s string) string {
	return html.EscapeString(s)
}

// HTMLUnescape 对HTML转义的字符串进行反转义
func HTMLUnescape(s string) string {
	return html.UnescapeString(s)
}

// XOREncrypt 使用异或运算对数据进行加密
func XOREncrypt(data []byte, key []byte) []byte {
	encrypted := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		encrypted[i] = data[i] ^ key[i%len(key)]
	}
	return encrypted
}

// XORDecrypt 使用异或运算对数据进行解密
func XORDecrypt(data []byte, key []byte) []byte {
	return XOREncrypt(data, key) // XOR解密和加密操作相同
}

// CaesarEncrypt 使用凯撒密码对字符串进行加密
func CaesarEncrypt(s string, shift int) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			result.WriteRune('A' + (r-'A'+rune(shift))%26)
		case r >= 'a' && r <= 'z':
			result.WriteRune('a' + (r-'a'+rune(shift))%26)
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// CaesarDecrypt 使用凯撒密码对字符串进行解密
func CaesarDecrypt(s string, shift int) string {
	return CaesarEncrypt(s, 26-shift)
}
