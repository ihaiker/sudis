package errors

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"error"`
	Message string `json:"message"`
}

func (err *Error) Error() string {
	return err.Code + ":" + err.Message
}

func New(format string, args ...interface{}) error {
	return NewError(iris.StatusInternalServerError, "InternalServerError", fmt.Sprintf(format, args...))
}

func NewError(status int, code string, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

var (
	ErrProgramNotFound = NewError(iris.StatusNotFound, "ErrProgramNotFound", "管托程序未发现")
	ErrProgramExists   = NewError(iris.StatusNotImplemented, "ErrProgramExists", "托管程序已经存在")

	ErrProgramIsRunning = NewError(iris.StatusNotFound, "ErrProgramIsRunning", "管托程序正在运行")

	ErrNodeIsEmpty = NewError(iris.StatusBadRequest, "NodeIsEmpty", "Node不能为空")
	ErrTagEmpty    = NewError(iris.StatusBadRequest, "TagIsEmpty", "Tag不能为空")
	ErrNameEmpty   = NewError(iris.StatusBadRequest, "NameIsEmpty", "Name不能为空")
	ErrClassEmpty  = NewError(iris.StatusBadRequest, "ClassEmpty", "class不能为空")
	ErrUser        = NewError(iris.StatusBadRequest, "UserError", "用户不存在或者密码不正确")

	ErrNotFound       = NewError(iris.StatusBadRequest, "NotFound", "未发现")
	ErrNotFoundConfig = NewError(iris.StatusBadRequest, "NotFoundConfig", "配置未发现")

	ErrNodeNotFound  = NewError(iris.StatusNotFound, "NodeNotFound", "节点未发现！")
	ErrNodeKeyExists = NewError(iris.StatusNotFound, "ErrNodeKeyExists", "节点主键已经存在，不可用！")

	ErrTimeout = NewError(iris.StatusInternalServerError, "ErrTimeout", "运行超时")

	ErrClientToken = NewError(iris.StatusUnauthorized, "ErrClientToken", "客户端Token错误")

	ErrToken = NewError(iris.StatusForbidden, "ErrToken", "错误的TOKEN")
)
