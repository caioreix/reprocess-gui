package utils

import (
	"fmt"
	"reflect"

	"github.com/golang-jwt/jwt/v5"
)

// NewJWT generates a new JWT token based on the provided object and key.
func NewJWT(obj any, key string) (string, error) {
	claims := structToClaims(obj)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ParseJWT parses the provided JWT token string using the given key and stores the result in the output.
func ParseJWT[T any](tokenString, key string, output T) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	err = claimsToStruct(claims, output)
	if err != nil {
		return err
	}

	return nil
}

func structToClaims(obj any) jwt.MapClaims {
	claims := make(jwt.MapClaims)

	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		value := val.Field(i).Interface()
		claims[tag] = value
	}

	return claims
}

func claimsToStruct[T any](claims jwt.MapClaims, output T) error {
	val := reflect.ValueOf(output).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}
		value, ok := claims[tag]
		if !ok {
			continue
		}

		fieldValue := val.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set value for field '%s'", field.Name)
		}

		claimValue := reflect.ValueOf(value)
		if !claimValue.CanConvert(field.Type) {
			return fmt.Errorf("cannot convert tag '%s' value '%s' to type '%s'", tag, value, field.Type)
		}
		convertedValue := claimValue.Convert(field.Type)

		if !convertedValue.Type().AssignableTo(field.Type) {
			return fmt.Errorf("cannot assign tag '%s' value '%s' to field '%s'", tag, value, field.Name)
		}

		fieldValue.Set(convertedValue)
	}

	return nil
}
