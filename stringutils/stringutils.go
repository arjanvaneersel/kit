package stringutils

import (
	"math/rand"
	"strings"
	"unicode"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
const (
	charBits = 6
	charMask = 1<<charBits - 1
	charMax  = 63 / charBits
)

func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), charMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), charMax
		}
		if j := int(cache & charMask); j < len(characters) {
			b[i] = characters[j]
			i--
		}
		cache >>= charBits
		remain--
	}
	return string(b)
}

func In(i bool, t string, args ...string) bool {
	var v1 string
	var v2 string

	// Loop over the provided args
	for _, a := range args {
		// Set the values to compare to lower case if i is true
		if i {
			v1 = strings.ToLower(t)
			v2 = strings.ToLower(a)
		} else {
			v1 = t
			v2 = a
		}

		// Compare the values
		if v1 == v2 {
			return true
		}
	}

	// Nothing found, return false
	return false
}

// replaceAll will replace any item under old with the value of new in the string
// Returns a string
func replaceAll(s string, new string, old ...string) string {
	for _, i := range old {
		s = strings.Replace(s, i, new, -1)
	}
	return s
}

// TitledString is a Stringer interface implementation, which replaces the characters - and _ by a space
// and capitalizes the first letter of the string
// Returns a string
type TitledString string

func (str TitledString) String() string {
	s := []rune(replaceAll(string(str), " ", "-", "_"))
	s[0] = unicode.ToUpper(s[0])
	return string(s)
}

//ReplaceBetween will replace occurances of t with st (StartTag) and et (EndTag)
func ReplaceBetween(s *string, t, st, et string) {
	var p int = len(t)
	var useSt bool = true
	str := *s

	for {
		i := strings.Index(str, t)
		if i == -1 {
			*s = str
			break
		}

		if useSt {
			str = str[:i] + st + str[i+p:]
		} else {
			str = str[:i] + et + str[i+p:]
		}
		useSt = !useSt
	}
}

func findBetween(s, st, et string, sp int) (int, int) {
	str := s[sp:]
	lst := len(st)
	i := strings.Index(str, st)
	if i == -1 {
		return -1, -1
	}

	j := strings.Index(str[i+lst:], et)
	if j == -1 {
		return i, -1
	}
	return i + lst, i + j - 1
}

func FindBetween(s, st, et string, sp int) string {
	str := s[sp:]
	start, end := findBetween(str, st, et, 0)
	if start == -1 || end == -1 {
		return ""
	}
	return str[start:end]
}
