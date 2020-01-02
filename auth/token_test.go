package auth

import "testing"

func TestCreateToken(t *testing.T) {
	// Create
	token := NewToken("1234")
	if token == nil {
		t.Errorf("expected token\n")
		return
	}

	// Encode
	tokenStr, err := token.Encode()
	if err != nil {
		t.Errorf("error not expected\ngot:%#v", err)
	}
	if len(tokenStr) < 10 {
		t.Errorf("wrong tokenStr length\n")
	}

	// Decode
	decodedToken, err := decodeToken(tokenStr)
	if err != nil {
		t.Errorf("error not expected\ngot:%#v", err)
		return
	}

	if token.ID.Hex() != decodedToken.ID.Hex() || token.UserID != decodedToken.UserID {
		t.Errorf("token and decodedToken are not equal\n-expected:%#v\n-actual:  %#v", token, decodedToken)
	}
}
