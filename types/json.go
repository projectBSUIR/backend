package types

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func ConvertMapToString(verdict any) (string, error) {
	ret, err := json.Marshal(verdict)
	return string(ret), err
}

func ConvertToMap(sverdict string) fiber.Map {
	var verdict fiber.Map
	_ = json.Unmarshal([]byte(sverdict), &verdict)
	return verdict
}
