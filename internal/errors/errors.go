package errors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// 自定義的 errors
var (
	ErrInvalidInput = &_error{Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest}

	ErrPageNotFound = &_error{Message: "Page Not Found.", Status: http.StatusNotFound}

	ErrResourceNotFound = &_error{Message: "Record Not Found.", Status: http.StatusNotFound}

	ErrResourceAlreadyExists = &_error{Message: "The specified resource already exists.", Status: http.StatusConflict}

	ErrInternalError = &_error{Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusBadRequest}

	ErrEmailNotFilledIn = &_error{Status: http.StatusBadRequest, Message: "email can't be empty",
		OtherLanguage: otherLanguage{Chinese: "信箱未填寫"}}
	ErrPhoneAreaCodeNotFilledIn = &_error{Status: http.StatusBadRequest, Message: "Phone number area code is not filled in",
		OtherLanguage: otherLanguage{Chinese: "電話號碼區碼未填寫"}}
	ErrPhoneNumberNotFilledIn = &_error{Status: http.StatusBadRequest, Message: "Phone number is not filled in",
		OtherLanguage: otherLanguage{Chinese: "電話號碼未填寫"}}
	ErrRegistrationTypeInvalidInput = &_error{Status: http.StatusBadRequest, Message: "Registration type is not supported",
		OtherLanguage: otherLanguage{Chinese: "註冊類型不支援"}}
	ErrNameNotFilledIn = &_error{Status: http.StatusBadRequest, Message: "Name has not been entered",
		OtherLanguage: otherLanguage{Chinese: "名字尚未輸入"}}
	ErrPasswordInvalidInput = &_error{Status: http.StatusBadRequest, Message: "Password is less than 8 characters",
		OtherLanguage: otherLanguage{Chinese: "密碼少於8個字元"}}
	ErrConfirmPasswordNotFilledIn = &_error{Status: http.StatusBadRequest, Message: "Confirm password cannot be empty",
		OtherLanguage: otherLanguage{Chinese: "確認密碼不可為空"}}
	ErrConfirmPasswordIncorrect = &_error{Status: http.StatusBadRequest, Message: "Confirm password is incorrect",
		OtherLanguage: otherLanguage{Chinese: "確認密碼不正確"}}

	ErrEmailAlreadyExists = &_error{Status: http.StatusConflict, Message: "Email being used",
		OtherLanguage: otherLanguage{Chinese: "Email已被使用"}}
	ErrPhoneAlreadyExists = &_error{Status: http.StatusConflict, Message: "Phone number being used",
		OtherLanguage: otherLanguage{Chinese: "電話號碼已被使用"}}
)

type _error struct {
	Status        int           `json:"status"`
	Message       string        `json:"message"`
	OtherLanguage otherLanguage `json:"-"`
}

type otherLanguage struct {
	Chinese string
}

// HTTPError ...
type HTTPError struct {
	ErrorMessage string `json:"error_message"`
}

func (e *_error) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Status))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// Is ...
func (e *_error) Is(target error) bool {
	causeErr := errors.Cause(target)
	tErr, ok := causeErr.(*_error)
	if !ok {
		return false
	}
	return e.Message == tErr.Message
}

// GetHTTPError ,,,
func GetHTTPError(c echo.Context, err *_error) HTTPError {
	msg := err.Message

	if strings.ContainsAny(c.Request().Header.Get("Accept-Language"), "ch") &&
		err.OtherLanguage.Chinese != "" {
		msg = err.OtherLanguage.Chinese
	}

	return HTTPError{
		ErrorMessage: msg,
	}
}

// NewWithMessage 抽換錯誤訊息
// 未定義的錯誤會被視為 ErrInternalError 類型
func NewWithMessage(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return WithStack(&_error{
			Status:  ErrInternalError.Status,
			Message: ErrInternalError.Message,
		})
	}
	err = &_error{
		Status:  _err.Status,
		Message: message,
	}
	var msg string
	for i := 0; i < len(args); i++ {
		msg += "%+v"
	}
	return Wrapf(err, msg, args...)
}

// WithErrors 使用訂好的errors code 與訊息,如果未定義message 顯示對應的http status描述
func WithErrors(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return WithStack(&_error{
			Status:  ErrInternalError.Status,
			Message: http.StatusText(ErrInternalError.Status),
		})
	}
	return WithStack(&_error{
		Status:  _err.Status,
		Message: _err.Message,
	})
}

// NewWithMessagef 抽換錯誤訊息
func NewWithMessagef(err error, format string, args ...interface{}) error {
	return NewWithMessage(err, fmt.Sprintf(format, args...))
}
