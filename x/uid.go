package x

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// GenerateUUID 生成唯一id
func GenerateUUID() string {
	u := uuid.Must(uuid.NewV4()).String()
	return strings.Replace(u, "-", "", -1)
}
