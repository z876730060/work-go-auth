package common

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("123456")

// ParseID 将字符串ID转换为uint类型，包含错误处理和边界检查
func ParseID(id string) (uint, error) {
	if id == "" {
		return 0, fmt.Errorf("id is empty")
	}

	// 使用更安全的转换方式，明确指定位数为64
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id format: %s", id)
	}

	// 检查ID是否超出uint范围
	if uid > math.MaxUint {
		return 0, fmt.Errorf("id too large: %s", id)
	}

	return uint(uid), nil
}

func RespOk(message string, data any, info any) map[string]any {
	return map[string]any{
		"code":    http.StatusOK,
		"message": message,
		"data":    data,
		"info":    info,
	}
}

func RespErr(message string, info any) map[string]any {
	return map[string]any{
		"code":    http.StatusInternalServerError,
		"message": message,
		"info":    info,
	}
}

type CompatibleClaims struct {
	UserID   uint     `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// 生成与 Java 兼容的 JWT
func GenerateCompatibleToken(userID uint, username string, roles []string) (string, error) {
	// 确保密钥长度
	secretKey := ensureKeyLength(JWTSecret)

	claims := CompatibleClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-app",                               // 必须与 Java 一致
			Subject:   strconv.FormatUint(uint64(userID), 10), // 必须设置
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        generateJWTID(), // 可选：JWT ID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 验证 Java 生成的 JWT
func ValidateJavaJWT(tokenString string) (*CompatibleClaims, error) {
	secretKey := ensureKeyLength(JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, &CompatibleClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CompatibleClaims); ok && token.Valid {
		// 验证必要字段
		if claims.Issuer != "my-app" {
			return nil, fmt.Errorf("签发者不匹配")
		}
		if claims.Subject == "" {
			return nil, fmt.Errorf("主题不能为空")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("令牌无效")
}

// 确保密钥长度合适
func ensureKeyLength(secret []byte) []byte {
	if len(secret) < 32 {
		// 填充到 32 字节
		newSecret := make([]byte, 32)
		copy(newSecret, secret)
		return newSecret
	}
	return secret
}

func generateJWTID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
