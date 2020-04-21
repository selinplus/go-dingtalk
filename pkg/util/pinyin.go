package util

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
	"unicode"
)

// 汉字转拼音
func TransToCharacter(s string) string {
	a := pinyin.NewArgs()
	tmp := make([]string, 0)
	s = strings.TrimSpace(s)
	beg := strings.Index(s, "（")
	end := strings.Index(s, "）")
	if beg > 0 {
		if end == len(s)-3 {
			s = s[:beg]
		} else {
			if beg == 0 {
				s = s[end+3:]
			} else {
				s1 := s[:beg+3]
				s2 := s[end+3:]
				s = s1 + s2
			}
		}
	}
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "\\", "")
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "+", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "~", "")
	s = strings.ReplaceAll(s, "%", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "·", "")
	s = strings.ReplaceAll(s, "α", "")
	s = strings.ReplaceAll(s, "β", "")
	s = strings.ReplaceAll(s, "γ", "")
	s = strings.ReplaceAll(s, "δ", "")
	length := len(s)
	for i := 0; i < length; {
		c := s[i : i+1]
		rs := []rune(c)
		if unicode.IsDigit(rs[0]) || unicode.IsLetter(rs[0]) {
			tmp = append(tmp, string(c))
			i++
			continue
		} else {
			if i+3 > length {
				return s
			}
			c = s[i : i+3]
			quan := pinyin.Pinyin(string(c), a)
			for _, z := range quan {
				tmp = append(tmp, z[0][:1])
			}
			i = i + 3
		}
	}
	res := strings.Join(tmp, "")
	return res
}
