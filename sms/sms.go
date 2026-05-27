// Package sms 提供统一的短信发送能力，支持阿里云、腾讯云等服务商。
package sms

// Provider 短信服务商接口。所有服务商需实现此接口。
type Provider interface {
	// Send 发送短信。
	Send(req *SendRequest) (*SendResponse, error)
	// ProviderName 返回服务商名称，如 "aliyun"、"tencent"。
	ProviderName() string

	// TemplateList 查询模板列表。
	TemplateList(req *TemplateListRequest) (*TemplateListResponse, error)
	// ApplyTemplate 申请短信模板。
	ApplyTemplate(req *ApplyTemplateRequest) (*ApplyTemplateResponse, error)

	// SignatureList 查询签名列表。
	SignatureList(req *SignatureListRequest) (*SignatureListResponse, error)
	// ApplySignature 申请短信签名。
	ApplySignature(req *ApplySignatureRequest) (*ApplySignatureResponse, error)

	// SendStatistics 查询发送统计。
	SendStatistics(req *StatisticsRequest) (*StatisticsResponse, error)
	// QueryRecords 查询发送记录。
	QueryRecords(req *QueryRecordRequest) (*QueryRecordResponse, error)
}

// SendRequest 短信发送请求。
type SendRequest struct {
	// PhoneNumbers 目标手机号，多个号码由服务商各自约定的分隔符分隔。
	PhoneNumbers string
	// TemplateID 模板 ID，在云平台控制台配置。
	TemplateID string
	// TemplateParams 模板参数，key 为模板占位符名称，value 为实际值。
	// 阿里云示例：{"code": "123456"} → 模板中 ${code} 被替换为 123456。
	// 腾讯云按 key 字母序排列 value 传入，请确保 key 字母序与模板占位符顺序一致。
	TemplateParams map[string]string
	// SignName 短信签名，需在云平台报备。
	SignName string
}

// SendResponse 短信发送结果。
type SendResponse struct {
	RequestID string // 云平台返回的请求 ID
	BizID     string // 发送回执 ID
	Code      string // 状态码，"OK" 表示成功
	Message   string // 状态描述
}

// SendCode 发送短信验证码的便捷方法。
// params 为模板参数，如 map[string]string{"code": "123456"}。
// SignName 使用 Provider 配置中的默认签名，如需覆盖请直接调用 Send。
func SendCode(p Provider, phone, templateID string, params map[string]string) (*SendResponse, error) {
	return p.Send(&SendRequest{
		PhoneNumbers:   phone,
		TemplateID:     templateID,
		TemplateParams: params,
	})
}
