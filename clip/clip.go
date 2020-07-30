package clip

func GetClipboard() ([]byte, error) {
	scrn, err := getClipboard()
	return scrn, err
}
