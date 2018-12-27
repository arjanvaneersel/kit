package i18n

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const TagName string = "translate"

type Translations map[string]string

type errTranslating struct {
	s string
	e string
}

func (e errTranslating) Error() string {
	return fmt.Sprintf("%s: %q", e.s, e.e)
}

func (e errTranslating) SetTranslationError(s error) errTranslating {
	e.e = s.Error()
	return e
}

func NewErrTranslating(s string) errTranslating {
	return errTranslating{s: s}
}

var (
	LANG           string
	ErrNotPointer  = errors.New("Output is not a pointer.")
	ErrNotStruct   = errors.New("Data is not a struct")
	ErrTranslating = NewErrTranslating("Error while translating")
	ErrWrongType   = errors.New("Tag \"translate\" used on a field which is not of type I18nString")
)

type I18nField struct {
	Default string
	L       string
	M       Translations
}

func (s *I18nField) fromString(str string) error {
	if s.M == nil {
		s.M = make(Translations)
	}

	p := strings.Split(str, ";;")
	for i, v := range p {
		if i == 0 {
			s.Default = v
			continue
		}

		l := strings.Split(v, "::")
		if len(l) != 2 {
			return fmt.Errorf("%d: malformed entry %q", i, v)
		}
		s.M[l[0]] = l[1]
	}
	return nil
}

func (s *I18nField) toString() (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s", s.Default))

	for k, v := range s.M {
		buffer.WriteString(fmt.Sprintf(";;%s::%s", strings.ToUpper(k), v))
	}

	return buffer.String(), nil
}

func (s *I18nField) GetTranslation(lang string) (string, error) {
	tr, ok := s.M[strings.ToLower(lang)]
	if !ok {
		return "", nil
	}
	return tr, nil
}

func (s *I18nField) SetTranslation(lang, txt string) {
	if s.M == nil {
		s.M = make(Translations)
	}
	s.M[strings.ToLower(lang)] = txt
}

func (s *I18nField) String() string {
	if s.L != "" {
		str, _ := s.GetTranslation(s.L)
		return str
	}
	if LANG != "" {
		str, _ := s.GetTranslation(LANG)
		return str
	}
	return s.Default
}

func (s *I18nField) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Default      string            `json:"default"`
		Translated   string            `json:"translated"`
		Translations map[string]string `json:"translations"`
	}{
		Default:      s.Default,
		Translated:   s.String(),
		Translations: s.M,
	})
}

func (s *I18nField) UnmarshallJSON(b []byte) error {
	decoded := new(struct {
		Default      string            `json:"default"`
		Translated   string            `json:"translated"`
		Translations map[string]string `json:"translations"`
	})

	err := json.Unmarshal(b, decoded)
	if err != nil {
		return err
	}
	s.Default = decoded.Default
	s.M = decoded.Translations
	return nil
}

func (s *I18nField) GetBSON() (interface{}, error) {
	res, err := s.toString()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *I18nField) SetBSON(raw bson.Raw) error {
	var i string
	err := raw.Unmarshal(&i)
	if err != nil {
		return err
	}
	err = s.fromString(i)
	if err != nil {
		return err
	}
	return nil
}

func NewI18nField(s string, t Translations) I18nField {
	return I18nField{Default: s, M: t}
}

func SetStructLanguage(data interface{}, lang string) error {
	var m I18nField
	tp := reflect.TypeOf(m) // Type of I18nString

	//rt := reflect.TypeOf(data) // Reflection type of the input data
	rv := reflect.ValueOf(data)                 // Reflection value of the input data
	if rv.Kind() != reflect.Ptr || rv.IsNil() { // Return an error if the reflection value is not a pointer or nil
		return ErrNotPointer
	}

	d := reflect.Indirect(rv)           // Get the elements of input data
	for i := 0; i < d.NumField(); i++ { // Loop over all fields
		d.Field(i)                                 // Get the current field
		if d.Field(i).Type().Name() != tp.Name() { // Move on if the field is not an I18nField
			continue
		}
		l := d.Field(i).FieldByName("L")
		if !l.CanSet() {
			return fmt.Errorf("Can't set field")
		}
		if l.Kind() != reflect.String {
			return fmt.Errorf("Kind is not a string")
		}
		l.SetString(lang) // Set the language
	}
	return nil
}
