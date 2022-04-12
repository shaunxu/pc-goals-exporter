package cookieparser

import (
	"bufio"
	"os"
	"strings"
)

type NCFCookieParser struct {
	CookieParserStringify
	Path string
}

func (p *NCFCookieParser) Parse() (Cookies, error) {
	f, err := os.Open(p.Path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	cookies := make(Cookies)
	for scanner.Scan() {
		line := scanner.Text()
		// ignore comment lines
		if strings.HasPrefix(line, "#") {
			continue
		}

		elems := strings.Split(line, "\t")
		// ignore lines without necessary elements
		if len(elems) < 7 {
			continue
		}
		cookies[elems[5]] = elems[6]
	}
	return cookies, nil
}
