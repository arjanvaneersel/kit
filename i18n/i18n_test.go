package i18n

import (
	"encoding/json"
	"fmt"
	"testing"
)

var tstr = NewI18nField("Hello", nil)
var translationTests = []struct {
	lang string
	txt  string
}{
	{"nl", "Hallo"},
	{"fr", "Bonjour"},
	{"bg", "Здравей"},
}

var tstStruct = struct {
	Name    I18nField
	Age     int
	Comment I18nField
}{
	Name:    NewI18nField("Mister Gopher", Translations{"bg": "Госпидин Гофер", "nl": "Meneer Gopher"}),
	Age:     40,
	Comment: NewI18nField("Go rules!", Translations{"bg": "Го царува!", "nl": "Go regeert!"}),
}

func TestI18nField_GetTranslation(t *testing.T) {
	for i, test := range translationTests {
		tstr.SetTranslation(test.lang, test.txt)
		res, err := tstr.GetTranslation(test.lang)
		if err != nil {
			t.Fatalf("[FAIL] Test %d: %q", i, err)
		}
		if res != test.txt {
			t.Fatalf("[FAIL] Test %d: Expected result for language %q to be %q, but got %q instead.", i, test.lang, test.txt, res)
		}
	}
}

func TestI18nField_StringWithLangVar(t *testing.T) {
	LANG = ""
	res := tstr.String()
	if res != "Hello" {
		t.Fatalf("[FAIL] Expected result for default language to be \"Hello\", but got %q instead.", res)
	}

	for i, test := range translationTests {
		LANG = test.lang
		res := tstr.String()
		if res != test.txt {
			t.Fatalf("[FAIL] Test %d: Expected result for language %q to be %q, but got %q instead.", i, test.lang, test.txt, res)
		}
	}
}

func TestI18nField_StringWithLangField(t *testing.T) {
	LANG = ""
	tstr.L = ""
	res := tstr.String()
	if res != "Hello" {
		t.Fatalf("[FAIL] Expected result for default language to be \"Hello\", but got %q instead.", res)
	}

	for i, test := range translationTests {
		tstr.L = test.lang
		res := tstr.String()
		if res != test.txt {
			t.Fatalf("[FAIL] Test %d: Expected result for language %q to be %q, but got %q instead.", i, test.lang, test.txt, res)
		}
	}
}

func TestSetStructLanguage(t *testing.T) {
	lang := "nl"
	err := SetStructLanguage(&tstStruct, lang)
	if err != nil {
		t.Fatalf("[FAIL] Error SetStructLanguage returned: %q", err)
	}
	if tstStruct.Name.L != lang {
		t.Fatalf("[FAIL] Expected l of Name to be %q, but got %q instead.", lang, tstStruct.Name.L)
	}

	if tstStruct.Comment.L != lang {
		t.Fatalf("[FAIL] Expected l of Comment to be %q, but got %q instead.", lang, tstStruct.Comment.L)
	}

	res := fmt.Sprintf("%s: %s", tstStruct.Name.String(), tstStruct.Comment.String())
	e := "Meneer Gopher: Go regeert!"
	if res != e {
		t.Fatalf("[FAIL] Expected %q, but got %q instead.", e, res)
	}
}

func TestI18nField_JSON(t *testing.T) {
	lang := "nl"
	err := SetStructLanguage(&tstStruct, lang)
	if err != nil {
		t.Fatalf("[FAIL] Error SetStructLanguage returned: %q", err)
	}

	res, err := json.Marshal(&tstStruct)
	if err != nil {
		t.Fatalf("[FAIL] Error json.Marshal returned: %q", err)
	}

	var res2 = struct {
		Name    I18nField
		Age     int
		Comment I18nField
	}{}

	err = json.Unmarshal(res, &res2)
	if err != nil {
		t.Fatalf("[FAIL] Error json.Unmarshal returned: %q", err)
	}

	if res2.Name.Default == tstStruct.Name.Default &&
		res2.Name.L == tstStruct.Name.L &&
		res2.Name.M["nl"] == tstStruct.Name.M["nl"] &&
		res2.Name.M["bg"] == tstStruct.Name.M["bg"] {
		t.Fatalf("[FAIL] Result doesn't match with original struct.\nGot: %+v.\nExpected: %+v", res2, tstStruct)
	}

	if res2.Comment.Default == tstStruct.Comment.Default &&
		res2.Comment.L == tstStruct.Comment.L &&
		res2.Comment.M["nl"] == tstStruct.Comment.M["nl"] &&
		res2.Comment.M["bg"] == tstStruct.Comment.M["bg"] {
		t.Fatalf("[FAIL] Result doesn't match with original struct.\nGot: %+v.\nExpected: %+v", res2, tstStruct)
	}
}

// Todo: Test bson (un)marshalling
// TOdo: Implement support for other databases
