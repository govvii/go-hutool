package desensitized

import (
	"strings"
)

// IDCardNum 身份证号码脱敏
// 保留前N位和后M位，其他用星号替换
func IDCardNum(idCard string, front, end int) string {
	if len(idCard) < front+end {
		return strings.Repeat("*", len(idCard))
	}
	return idCard[:front] + strings.Repeat("*", len(idCard)-front-end) + idCard[len(idCard)-end:]
}

// MobilePhone 手机号脱敏
// 保留前3位和后4位，中间用4个星号替换
func MobilePhone(phone string) string {
	if len(phone) != 11 {
		return strings.Repeat("*", len(phone))
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// Password 密码脱敏
// 将所有字符替换为星号
func Password(password string) string {
	return strings.Repeat("*", len(password))
}

// ChineseName 中文姓名脱敏
// 保留姓氏，其他用星号替换
func ChineseName(name string) string {
	if len(name) <= 1 {
		return name
	}
	return name[:1] + strings.Repeat("*", len(name)-1)
}

// Email 电子邮件脱敏
// 隐藏用户名中间部分
func Email(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return strings.Repeat("*", len(email))
	}
	username, domain := parts[0], parts[1]
	if len(username) <= 2 {
		return username + "@" + domain
	}
	return username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:] + "@" + domain
}

// BankCard 银行卡号脱敏
// 保留前6位和后4位
func BankCard(cardNo string) string {
	if len(cardNo) <= 10 {
		return strings.Repeat("*", len(cardNo))
	}
	return cardNo[:6] + strings.Repeat("*", len(cardNo)-10) + cardNo[len(cardNo)-4:]
}

// Address 地址脱敏
// 保留前6个字符和后6个字符
func Address(address string) string {
	if len(address) <= 12 {
		return strings.Repeat("*", len(address))
	}
	return address[:6] + strings.Repeat("*", len(address)-12) + address[len(address)-6:]
}

// LicensePlate 车牌号脱敏
// 保留前两位和最后一位
func LicensePlate(plate string) string {
	if len(plate) <= 3 {
		return strings.Repeat("*", len(plate))
	}
	return plate[:2] + strings.Repeat("*", len(plate)-3) + plate[len(plate)-1:]
}

// Landline 座机号脱敏
// 保留区号和后4位
func Landline(phone string) string {
	parts := strings.Split(phone, "-")
	if len(parts) != 2 {
		return strings.Repeat("*", len(phone))
	}
	areaCode, number := parts[0], parts[1]
	if len(number) <= 4 {
		return areaCode + "-" + strings.Repeat("*", len(number))
	}
	return areaCode + "-" + strings.Repeat("*", len(number)-4) + number[len(number)-4:]
}
