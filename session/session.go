package session

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

type SessionValues map[string]interface{}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func CreateSession(w http.ResponseWriter, name string, values *SessionValues) {
	if encoded, err := cookieHandler.Encode(name, *values); err == nil {
		c := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
	}
}

func DestroySession(w http.ResponseWriter, name string) {
	c := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}

func GetValues(name string, r *http.Request) (SessionValues, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}

	sValues := make(SessionValues)
	if err = cookieHandler.Decode(name, c.Value, &sValues); err != nil {
		return nil, err
	}
	return sValues, nil
}

func Get

//Todo: Tests, comments and documentation
