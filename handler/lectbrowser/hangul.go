package lectbrowser

import "bytes"

const (
	firstConsonant = 'ㄱ'
	lastConsonant  = 'ㅎ'
	firstCharacter = '가'
	lastCharacter  = '힣'
	footer         = 'ㅣ' //588
)

var consonants = [...]rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ',
	'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

func HangulsToConsonants(str string) string {
	var buf bytes.Buffer
	for _, ch := range str {
		buf.WriteRune(HangulToConsonant(ch))
	}
	return buf.String()
}

func HangulToConsonant(ch rune) rune {
	if IsHangul(ch) {
		return consonants[(ch-firstCharacter)/588]
	}
	return ch
}

func IsConsonant(ch rune) bool {
	return firstConsonant <= ch && ch <= lastConsonant
}
func IsHangul(ch rune) bool {
	return firstCharacter <= ch && ch <= lastCharacter
}
