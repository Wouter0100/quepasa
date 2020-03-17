package controllers

// import (
// 	"testing"
// )

// func TestParseJSONBody(t *testing.T) {
// 	t.Run("Valid POST", func(t *testing.T) {
// 		event := map[string]interface{}{
// 			"body": `{"number": "15555555555", "message": "This is a message"}`,
// 		}
// 		_, err := parseJSONBody(event)
// 		if err != nil {
// 			t.Error(err.Error())
// 		}
// 	})

// 	t.Run("Invalid POST", func(t *testing.T) {
// 		event := map[string]interface{}{
// 			"body": `{"message": A message"}`,
// 		}
// 		_, err := parseJSONBody(event)
// 		if err.Error() != "invalid character 'A' looking for beginning of value" {
// 			t.Error("expected invalid JSON string")
// 		}
// 	})
// }
