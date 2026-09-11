package validation

import "testing"

func TestValidEmail(t *testing.T) {
	for _, test := range []struct {
		value string
		want  bool
	}{{"prueba@gmail.com", true}, {"sin-arroba.com", false}, {"a@b", false}, {"nombre <a@b.com>", false}} {
		if got := ValidEmail(test.value); got != test.want {
			t.Errorf("ValidEmail(%q) = %v; se esperaba %v", test.value, got, test.want)
		}
	}
}

func TestValidPhone(t *testing.T) {
	for _, test := range []struct {
		value string
		want  bool
	}{{"5512345678", true}, {"+525512345678", true}, {"123", false}, {"05512345678", false}, {"55abc45678", false}} {
		if got := ValidPhone(test.value); got != test.want {
			t.Errorf("ValidPhone(%q) = %v; se esperaba %v", test.value, got, test.want)
		}
	}
}

func TestValidPassword(t *testing.T) {
	for _, test := range []struct {
		value string
		want  bool
	}{{"Abc1@x", true}, {"Abcdef1@xy12", true}, {"Ab1@x", false}, {"abcdefghijkl", false}, {"ABCDEF1@", false}, {"abcdef1@", false}, {"Abcdefg@", false}, {"Abcdef1!", false}, {"Abcdef1@xy123", false}} {
		if got := ValidPassword(test.value); got != test.want {
			t.Errorf("ValidPassword(%q) = %v; se esperaba %v", test.value, got, test.want)
		}
	}
}
