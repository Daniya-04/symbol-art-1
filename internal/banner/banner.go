package banner

func BannerCheck(input string) map[string][][]string {

}

func Validate(sl []rune) bool {

	for _, v := range sl {
		if v <= 32 || v >= 126 {
			return false
		}
	}
	return true
}
