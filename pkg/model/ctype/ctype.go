package ctype

// VerifyStatus 驗證狀態
type VerifyStatus int8

const (
	// VerifyInit 待寄送
	VerifyInit VerifyStatus = iota + 1
	// VerifySend 已寄送
	VerifySend
	// VerifySuccess 驗證成功
	VerifySuccess
	// VerifyFail 驗證失敗
	VerifyFail
)
