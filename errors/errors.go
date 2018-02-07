package errors

import "github.com/gin-gonic/gin"

func newError(code, text string) gin.H {
	return gin.H{"error": code, "text": text}
}

func IncorrectUserPwd() gin.H {
	return newError("001", "Incorrect username or password")
}

func CanntFindToken() gin.H {
	return newError("011", "Cann't find token!")
}

func NotEvenToken() gin.H {
	return newError("012", "That's not even a token!")
}

func ExpiredOrNotActiveToken() gin.H {
	return newError("013", "Token is either expired or not active yet!")
}

func InvalidUser() gin.H {
	return newError("014", " Invalid user name!")
}

func CanntHandleToken() gin.H {
	return newError("019", "Couldn't handle this token!")
}

func CanntGetData() gin.H {
	return newError("021", "Couldn't get data!")
}

func CanntGetParam(name string) gin.H {
	return newError("022", "Couldn't get param! "+name)
}

func CanntSaveData() gin.H {
	return newError("023", "Couldn't save data!")
}
