package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 自定義的 errors, 我還要定義更多 error \(ˋ皿ˊ)/  <---怪叔叔 ಠ_ಠ
var (
	// ErrBadRequest                       =  &_error{Code: "400000", Message: http.StatusText(http.StatusBadRequest), Status: http.StatusBadRequest}
	ErrInvalidInput  = &_error{Code: "400001", Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidAmount = &_error{Code: "400002", Message: "This amount is not allow ", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrSuccessRate   = &_error{Code: "400003", Message: "This SuccessRate is not allow ", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	// 興業銀行序列號錯誤
	ErrInvalidSerialCode          = &_error{Code: "400007", Message: "please entry the correct number", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidQueryParameterValue = &_error{Code: "400009", Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidHeaderValue         = &_error{Code: "400004", Message: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrMissingRequiredHeader      = &_error{Code: "400017", Message: "A required HTTP header was not specified.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrOutOfRangeInput            = &_error{Code: "400020", Message: "something out of range", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidAppVersion          = &_error{Code: "400030", Message: "Check app version from x-app-version, and the version is invalid", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInternalDataNotSync        = &_error{Code: "400041", Message: "The specified data not sync in system, please wait a moment.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrFishGeneralError           = &_error{Code: "400042", Message: "fish error", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrNotMatchSetting            = &_error{Code: "400087", Message: "The specified data not match setting, please adjust your inputs.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}

	ErrUnauthorized                = &_error{Code: "401001", Message: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrInvalidAuthenticationInfo   = &_error{Code: "401001", Message: "The authentication information was not provided in the correct format. Verify the value of Authorization header.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrUsernameOrPasswordIncorrect = &_error{Code: "401002", Message: "Username or Password is incorrect.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}

	// 支付寶爬蟲錯誤
	ErrCrawlSecurituCheck = &_error{Code: "401003", Message: "alipay crawler need pass security check", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrCrawlNeedLogin     = &_error{Code: "401004", Message: "alipay crawler need login again", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}

	// ErrForbidden                                   =  &_error{Code: "403000", Message: http.StatusText(http.StatusForbidden), Status: http.StatusForbidden}
	ErrAccountIsDisabled                           = &_error{Code: "403001", Message: "The specified account is disabled.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrAuthenticationFailed                        = &_error{Code: "403002", Message: "Server failed to authenticate the request. Make sure the value of the Authorization header is formed correctly including the signature.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrNotAllowed                                  = &_error{Code: "403003", Message: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpExpired                                  = &_error{Code: "403004", Message: "OTP is expired.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrInsufficientAccountPermissionsWithOperation = &_error{Code: "403005", Message: "The account being accessed does not have sufficient permissions to execute this operation.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpRequired                                 = &_error{Code: "403007", Message: "OTP Binding is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpAuthorizationRequired                    = &_error{Code: "403008", Message: "Two-factor authorization is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpIncorrect                                = &_error{Code: "403009", Message: "OTP is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrResetPasswordRequired                       = &_error{Code: "403010", Message: "Reset password is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrResetOTPAccountNameIncorrect                = &_error{Code: "403011", Message: "Reset otp failed,requested otp name is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrSignIncorrect                               = &_error{Code: "403012", Message: "verify sign failed", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrResetOTPAccountEmailIncorrect               = &_error{Code: "403013", Message: "Reset otp failed,requested otp email is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrRiskControlBlockCreateOrder                 = &_error{Code: "403014", Message: "风控禁止建单", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}

	// ErrNotFound         =  &_error{Code: "404000", Message: http.StatusText(http.StatusNotFound), Status: http.StatusNotFound}
	ErrResourceNotFound    = &_error{Code: "404001", Message: "The specified resource does not exist.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrAccountNotFound     = &_error{Code: "404002", Message: "cant find any account.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrPageNotFound        = &_error{Code: "404003", Message: "Page Not Found.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrOrderNotFound       = &_error{Code: "404004", Message: "The specified order not found", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrAccountNotAvailable = &_error{Code: "404012", Message: "account is not available", Status: http.StatusNotFound, GRPCCode: codes.NotFound}

	ErrMethodNotAllowed = &_error{Code: "405001", Message: "Server has received and recognized the request, but has rejected the specific HTTP method it’s using.", Status: http.StatusMethodNotAllowed, GRPCCode: codes.Unavailable}

	ErrRequestTime = &_error{Code: "408001", Message: "request time out", Status: http.StatusRequestTimeout, GRPCCode: codes.DeadlineExceeded}

	ErrConflict                 = &_error{Code: "409000", Message: http.StatusText(http.StatusConflict), Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrAccountAlreadyExists     = &_error{Code: "409001", Message: "The specified account already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrAccountBeingCreated      = &_error{Code: "409002", Message: "The specified account is in the process of being created.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrResourceAlreadyExists    = &_error{Code: "409004", Message: "The specified resource already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrPhoneVerifiedTimeout     = &_error{Code: "409007", Message: "sms verify time out", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrCreateRechargeOrderCDing = &_error{Code: "409008", Message: "60秒内只能建立一次充值申请", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}

	ErrInternalServerError = &_error{Code: "500000", Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
	ErrInternalError       = &_error{Code: "500001", Message: "The server encountered an internal error. Please retry the request.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}

	// Internal usage
	ErrOrderNoChange  = &_error{Code: "500002", Message: "Order status No change", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
	ErrGAMSettleError = &_error{Code: "500003", Message: "GAM settle failed", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)

type _error struct {
	Status   int                    `json:"status"`
	Code     string                 `json:"code"`
	GRPCCode codes.Code             `json:"grpccode"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details"`
}

// HttpError ...
type HttpError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func (e *_error) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(e.Code)
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
	return e.Code == tErr.Code
}

// GetHTTPError ,,,
func GetHTTPError(err *_error) HttpError {
	return HttpError{
		Message: err.Message,
		Code:    err.Code,
		Details: err.Details,
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
			Status:   ErrInternalError.Status,
			Code:     ErrInternalError.Code,
			Message:  ErrInternalError.Message,
			GRPCCode: ErrInternalError.GRPCCode,
		})
	}
	err = &_error{
		Status:   _err.Status,
		Code:     _err.Code,
		Message:  message,
		GRPCCode: _err.GRPCCode,
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
			Code:    ErrInternalError.Code,
			Message: http.StatusText(ErrInternalError.Status),
		})
	}
	return WithStack(&_error{
		Status:  _err.Status,
		Code:    _err.Code,
		Message: _err.Message,
	})
}

// SetDetails set details as you wish =)
func (e *_error) SetDetails(details map[string]interface{}) {
	e.Details = details
	return
}

// CompareErrorCode 比較兩個錯誤代碼是否一致
func CompareErrorCode(errA error, errB error) bool {
	var aErr, bErr *_error
	if err, exists := errors.Cause(errA).(*_error); exists {
		aErr = err
	}
	if err, exists := errors.Cause(errB).(*_error); exists {
		bErr = err
	}
	if aErr.Code == bErr.Code {
		return true
	}
	return false
}

//ConvertProtoErr Convert _error to grpc error
func ConvertProtoErr(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return status.Error(ErrInternalError.GRPCCode, err.Error())
	}
	b, _ := json.Marshal(_err)
	return status.Error(_err.GRPCCode, string(b))
}

//ConvertHttpErr Convert  grpc error to _error
func ConvertHttpErr(err error) error {
	if err == nil {
		return nil
	}
	s := status.Convert(err)
	if s == nil {
		return ErrInternalError
	}
	interErr := _error{}
	jerr := json.Unmarshal([]byte(s.Message()), &interErr)
	if jerr != nil {
		return switchCode(s)
	}
	return WithStack(&interErr)
}
func switchCode(s *status.Status) error {
	httperr := ErrInternalError
	switch s.Code() {
	case Unknown:
		httperr = ErrInternalError
	case InvalidArgument:
		httperr = ErrInvalidInput
	case NotFound:
		httperr = ErrResourceNotFound
	case AlreadyExists:
		httperr = ErrResourceAlreadyExists
	case PermissionDenied:
		httperr = ErrNotAllowed
	case Unauthenticated:
		httperr = ErrUnauthorized
	case OutOfRange:
		httperr = ErrInvalidInput
	case Internal:
		httperr = ErrInternalError
	case DataLoss:
		httperr = ErrInternalError
	}
	httperr.Message = s.Message()
	return WithStack(httperr)
}

// ConvertPostgresError convert postgres error
func ConvertPostgresError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResourceNotFound
	}

	return errors.WithMessage(ErrInternalError, err.Error())

}

// ConvertMySQLError convert mysql error
func ConvertMySQLError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrResourceNotFound
	}

	return errors.WithMessage(ErrInternalError, err.Error())
}

// NewWithMessagef 抽換錯誤訊息
func NewWithMessagef(err error, format string, args ...interface{}) error {
	return NewWithMessage(err, fmt.Sprintf(format, args...))
}

// GetCodeWithErrors 使用訂好的errors code 與訊息,如果未定義message 顯示對應的http status描述
func GetCodeWithErrors(err error) (string, string) {

	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return ErrInternalError.Code, ErrInternalError.Message
	}
	return _err.Code, _err.Message
}

// HTTPConvertToError 將 http 的 response body convert to _error
func HTTPConvertToError(b []byte) error {
	interErr := _error{}
	jErr := json.Unmarshal(b, &interErr)
	if jErr != nil {
		return ErrInternalError
	}
	return WithStack(&interErr)
}
