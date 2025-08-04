package helper

func SetPropertyRequest(status int) string {
	if status >= 400 {
		return  "error"
	} else {
		return  "message"
	}
}