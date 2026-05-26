package sms

import "errors"

var (
	// ErrMissingPhoneNumbers 未提供手机号。
	ErrMissingPhoneNumbers = errors.New("sms: phone numbers is required")
	// ErrMissingTemplateID 未提供模板 ID。
	ErrMissingTemplateID = errors.New("sms: template id is required")
	// ErrMissingSignName 未提供短信签名。
	ErrMissingSignName = errors.New("sms: sign name is required")
)
