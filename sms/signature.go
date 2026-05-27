package sms

// SignatureItem 签名信息。
type SignatureItem struct {
	SignName   string // 签名名称
	Status     string // 审核状态: "approved" 通过, "pending" 审核中, "rejected" 未通过
	Reason     string // 审核驳回原因
	CreateTime string // 创建时间
}

// SignatureListRequest 查询签名列表请求。
type SignatureListRequest struct {
	Page     int // 页码，从 1 开始
	PageSize int // 每页数量
}

// SignatureListResponse 查询签名列表响应。
type SignatureListResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TotalCount int
	Signatures []SignatureItem
}

// ApplySignatureRequest 申请短信签名请求。
type ApplySignatureRequest struct {
	SignName    string // 签名名称
	Remark      string // 申请说明
	SignSource  int    // 签名来源: 0=企事业单位全称或简称, 1=工信部备案网站域名
	ProofBase64 string // 资质证明图片 base64 编码
	ProofSuffix string // 资质证明文件后缀: jpg, png, pdf 等
}

// ApplySignatureResponse 申请短信签名响应。
type ApplySignatureResponse struct {
	RequestID string
	Code      string // "OK" 表示成功
	Message   string
	SignName  string // 签名名称
}

// ModifySignatureRequest 修改短信签名请求。
// 阿里云通过 SignName 定位要修改的签名；腾讯云需要传 SignID。
type ModifySignatureRequest struct {
	SignID      uint64 // 签名 ID（腾讯云必填）
	SignName    string // 签名名称
	Remark      string // 修改说明
	SignSource  int    // 签名来源: 0=企事业单位, 1=工信部备案域名
	ProofBase64 string // 资质证明图片 base64
	ProofSuffix string // 文件后缀: jpg/png/pdf
}

// ModifySignatureResponse 修改短信签名响应。
type ModifySignatureResponse struct {
	RequestID string
	Code      string // "OK" 表示成功
	Message   string
	SignName  string // 签名名称
}

// DeleteSignatureRequest 删除短信签名请求。
// 阿里云通过 SignName 定位要删除的签名；腾讯云需要传 SignID。
type DeleteSignatureRequest struct {
	SignID   uint64 // 签名 ID（腾讯云必填）
	SignName string // 签名名称（阿里云必填）
}

// DeleteSignatureResponse 删除短信签名响应。
type DeleteSignatureResponse struct {
	RequestID string
	Code      string // "OK" 表示成功
	Message   string
}
