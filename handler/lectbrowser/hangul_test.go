package lectbrowser

import "testing"

func TestHangulsToConsonants(t *testing.T) {
	inputs := []string{
		"C프로그래밍",
		"운영체제 개론",
	}
	wants := []string{
		"Cㅍㄹㄱㄹㅁ",
		"ㅇㅇㅊㅈ ㄱㄹ",
	}

	for i, v := range inputs {
		result := HangulsToConsonants(v)
		if result != wants[i] {
			t.Errorf("f(%v)=%v, want=%v", v, result, wants[i])
		}

	}

}

func TestHangulToConsonant(t *testing.T) {
	inputs := []rune("가나다라마바사아자차타카파하ㅂㅈㄷㄱㅅㅁㄴㅇㄹㅎㅋㅌㅊㅍ@#$%^&*(ertyuiopasdfghjklzxcvbnm")
	wants := []rune("ㄱㄴㄷㄹㅁㅂㅅㅇㅈㅊㅌㅋㅍㅎㅂㅈㄷㄱㅅㅁㄴㅇㄹㅎㅋㅌㅊㅍ@#$%^&*(ertyuiopasdfghjklzxcvbnm")

	for i, v := range inputs {
		result := HangulToConsonant(v)
		if result != wants[i] {
			t.Errorf("f(%c) = %c, want = %c", v, result, wants[i])
		}
	}
}

func TestIsConsonant(t *testing.T) {
	consonantInputs := "ㅂㅈㄷㄱㅅㅁㄴㅇㄹㅋㅎㅌㅊㅍㅃㅉㄸㄲㅆ"
	otherInputs := "다람쥐 헌 책바퀴asdfqwer!@#$%^&*"

	for _, c := range consonantInputs {
		if !IsConsonant(c) {
			t.Errorf("%c is consonant but was returned false", c)
		}
	}

	for _, o := range otherInputs {
		if IsConsonant(o) {
			t.Errorf("%c is not consonant but was returned true", o)
		}
	}
}

func TestIsHangul(t *testing.T) {
	hanguls := "가나다라마바사아자차카타파하힣"
	others := "ㅂㅈㄷㄱㅅㅁㄴㅇㄹㅋㅎㅌㅊㅍㅃㅉㄸㄲㅆqwertyuiopasdfghjklzxcvbnm,!@#$%^&*()_"

	for _, h := range hanguls {
		if !IsHangul(h) {
			t.Errorf("%c is Hangul but was returned false", h)
		}
	}

	for _, o := range others {
		if IsHangul(o) {
			t.Errorf("%c is not Hangul but was returned true", o)
		}
	}

}
