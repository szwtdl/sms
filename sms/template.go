package sms

// TemplateItem 模板信息。
type TemplateItem struct {
	TemplateID      string // 模板 ID
	TemplateName    string // 模板名称
	TemplateContent string // 模板内容
	Status          string // 审核状态: "approved" 通过, "pending" 审核中, "rejected" 未通过
	Reason          string // 审核驳回原因
	CreateTime      string // 创建时间
}

// TemplateListRequest 查询模板列表请求。
type TemplateListRequest struct {
	Page     int // 页码，从 1 开始
	PageSize int // 每页数量
}

// TemplateListResponse 查询模板列表响应。
type TemplateListResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TotalCount int
	Templates  []TemplateItem
}

// ApplyTemplateRequest 申请短信模板请求。
type ApplyTemplateRequest struct {
	TemplateName    string // 模板名称
	TemplateContent string // 模板内容，如 "您的验证码为{1}"
	Remark          string // 申请说明
	TemplateType    int    // 模板类型: 0=验证码, 1=通知, 2=推广
}

// ApplyTemplateResponse 申请短信模板响应。
type ApplyTemplateResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TemplateID string // 平台分配的模板 ID
}

// ModifyTemplateRequest 修改短信模板请求。
type ModifyTemplateRequest struct {
	TemplateID      string // 要修改的模板 ID
	TemplateName    string // 新的模板名称
	TemplateContent string // 新的模板内容
	TemplateType    int    // 模板类型: 0=验证码, 1=通知, 2=推广
	Remark          string // 修改说明
}

// ModifyTemplateResponse 修改短信模板响应。
type ModifyTemplateResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TemplateID string // 模板 ID
}

// DeleteTemplateRequest 删除短信模板请求。
type DeleteTemplateRequest struct {
	TemplateID string // 模板 ID
}

// DeleteTemplateResponse 删除短信模板响应。
type DeleteTemplateResponse struct {
	RequestID string
	Code      string // "OK" 表示成功
	Message   string
}
