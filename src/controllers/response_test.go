package controllers

// import (
// 	"errors"
// 	"math"
// 	"testing"
// )

/*
func TestRespondSuccess(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		valid := map[string]interface{}{
			"number":  "15555555555",
			"message": "This is a message",
		}
		res, _ := RespondSuccess(valid)
		if res.StatusCode != 200 {
			t.Errorf("expected status code 200, got %d", res.StatusCode)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalid := math.Inf(1)
		res, err := RespondSuccess(invalid)
		if res.StatusCode != 500 {
			t.Errorf("expected status code 500, got %d", res.StatusCode)
		} else if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRespondBadRequest(t *testing.T) {
	message := "invalid POST body"
	res, err := RespondBadRequest(errors.New(message))
	if res.StatusCode != 400 {
		t.Errorf("expected status code 400, got %d", res.StatusCode)
	} else if err.Error() != message {
		t.Errorf("expected error '%s', got '%s'", message, err.Error())
	}
}

func TestRespondUnauthorized(t *testing.T) {
	message := "secret is incorrect"
	res, err := RespondUnauthorized(errors.New(message))
	if res.StatusCode != 401 {
		t.Errorf("expected status code 401, got %d", res.StatusCode)
	} else if err.Error() != message {
		t.Errorf("expected error '%s', got '%s'", message, err.Error())
	}
}

func TestRespondNotFound(t *testing.T) {
	message := "user does not exist"
	res, err := RespondNotFound(errors.New(message))
	if res.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", res.StatusCode)
	} else if err.Error() != message {
		t.Errorf("expected error '%s', got '%s'", message, err.Error())
	}
}

func TestRespondServerError(t *testing.T) {
	message := "connection failure"
	res, err := RespondServerError(errors.New(message))
	if res.StatusCode != 500 {
		t.Errorf("expected status code 500, got %d", res.StatusCode)
	} else if err.Error() != message {
		t.Errorf("expected error '%s', got '%s'", message, err.Error())
	}
}
*/
