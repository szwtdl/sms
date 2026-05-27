package sms

// StatisticsRequest 发送统计请求。
type StatisticsRequest struct {
	StartDate string // 开始日期 yyyyMMdd，如 "20240501"
	EndDate   string // 结束日期 yyyyMMdd
	Type      int    // 短信类型，1国内，2国际/港澳台
	SignName  string // 可选，按签名过滤
	Page      int    // 页码，从 1 开始
	PageSize  int    // 每页数量
}

// StatisticsResponse 发送统计响应。
type StatisticsResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TotalSent  int64 // 发送总数
	SuccessCnt int64 // 成功数
	FailCnt    int64 // 失败数
}

// QueryRecordRequest 查询发送记录请求。
type QueryRecordRequest struct {
	PhoneNumber string // 手机号
	SendDate    string // 发送日期 yyyyMMdd
	Page        int    // 页码，从 1 开始
	PageSize    int    // 每页数量
}

// SendRecordItem 单条发送记录。
type SendRecordItem struct {
	PhoneNumber  string // 手机号
	SendDate     string // 发送时间
	ReceiveDate  string // 接收时间
	TemplateCode string // 模板码
	Content      string // 短信内容
	Status       string // "success" 成功, "failed" 失败, "pending" 发送中
	ErrCode      string // 错误码
}

// QueryRecordResponse 查询发送记录响应。
type QueryRecordResponse struct {
	RequestID  string
	Code       string // "OK" 表示成功
	Message    string
	TotalCount int
	Records    []SendRecordItem
}
