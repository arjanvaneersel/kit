package jwt

import "testing"

func TestJwt(t *testing.T) {
	expected := map[string]interface{}{"foo": "bar"}

	token, err := CreateToken(expected)
	if err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}

	got, err := ParseToken(token)
	if err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}
	if got["foo"] != expected["foo"] {
		t.Errorf("received invalid data")
	}
}
