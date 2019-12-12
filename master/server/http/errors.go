package http

import "github.com/kataras/iris"

type Error struct {
	Status  int    `json:"-"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewError(status int, err string, message string) *Error {
	return &Error{Status: status, Error: err, Message: message}
}

func NewErrorMessage(status int, message string) *Error {
	return &Error{Status: status, Error: "BadRequest", Message: message}
}

func Assert(isTrue bool, err *Error) {
	if !isTrue {
		panic(err)
	}
}

func AssertErr(err error) {
	if err != nil {
		panic(&Error{Status: iris.StatusInternalServerError, Error: "InternalServerError", Message: err.Error()})
	}
}

var ErrNodeIsEmpty = NewError(iris.StatusBadRequest, "NodeIsEmpty", "Node不能为空")
var ErrTagEmpty = NewError(iris.StatusBadRequest, "TagIsEmpty", "Tag不能为空")
var ErrNameEmpty = NewError(iris.StatusBadRequest, "NameIsEmpty", "Name不能为空")
var ErrClassEmpty = NewError(iris.StatusBadRequest, "ClassEmpty", "class不能为空")
var ErrUser = NewError(iris.StatusBadRequest, "UserError", "用户不存在或者密码不正确")
var ErrNotFound = NewError(iris.StatusBadRequest, "NotFound", "程序未发现")
