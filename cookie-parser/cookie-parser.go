package cookieparser

import (
	"fmt"
	"strings"
)

type Cookies map[string]string

type CookieParser interface {
	Parse() (Cookies, error)
	String(c Cookies) string
}

type CookieParserStringify struct {
}

func (p CookieParserStringify) String(cookies Cookies) string {
	builder := strings.Builder{}
	for name, value := range cookies {
		builder.WriteString(fmt.Sprintf("%s=%s; ", name, value))
	}
	return builder.String()
}
