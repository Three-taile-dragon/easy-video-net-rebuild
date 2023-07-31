package jwts

import "testing"

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzU0MzY4NDMsInRva2VuIjoiMTAwNSJ9.YPdF87iUfHVcar4vq0ryAME0mvaaICUMBYIIPIo0Fls"
	ParseToken(tokenString, "dragon")
}
