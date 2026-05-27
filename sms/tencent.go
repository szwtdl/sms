package sms

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	smsapi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// TencentConfig 腾讯云短信配置。
type TencentConfig struct {
	SecretID  string
	SecretKey string
	AppID     string // 短信应用 ID，在控制台 应用管理 中查看
	SignName  string // 默认短信签名，Send 时 SignName 为空则使用此值
}

type tencentProvider struct {
	client   *smsapi.Client
	appID    string
	signName string
}

// NewTencent 创建腾讯云短信服务。
func NewTencent(cfg TencentConfig) (Provider, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := smsapi.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, err
	}
	return &tencentProvider{client: client, appID: cfg.AppID, signName: cfg.SignName}, nil
}

func (t *tencentProvider) ProviderName() string { return "tencent" }

func (t *tencentProvider) Send(req *SendRequest) (*SendResponse, error) {
	if req.SignName == "" {
		req.SignName = t.signName
	}
	if err := validate(req); err != nil {
		return nil, err
	}

	apiReq := smsapi.NewSendSmsRequest()
	apiReq.PhoneNumberSet = strPtrs(splitPhones(req.PhoneNumbers))
	apiReq.SmsSdkAppId = &t.appID
	apiReq.SignName = &req.SignName
	apiReq.TemplateId = &req.TemplateID
	apiReq.TemplateParamSet = mapValuesSorted(req.TemplateParams)

	apiResp, err := t.client.SendSms(apiReq)
	if err != nil {
		return nil, err
	}

	resp := &SendResponse{
		RequestID: *apiResp.Response.RequestId,
	}
	if len(apiResp.Response.SendStatusSet) > 0 {
		st := apiResp.Response.SendStatusSet[0]
		resp.Code = *st.Code
		resp.Message = *st.Message
	}
	return resp, nil
}

// mapValuesSorted 按 key 字母序返回 map 的 value 指针切片，
// 确保腾讯云模板参数顺序可预期。
func mapValuesSorted(m map[string]string) []*string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ptrs := make([]*string, len(keys))
	for i, k := range keys {
		ptrs[i] = &[]string{m[k]}[0] // 取 value 的指针
	}
	return ptrs
}

func splitPhones(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func (t *tencentProvider) TemplateList(req *TemplateListRequest) (*TemplateListResponse, error) {
	apiReq := smsapi.NewDescribeSmsTemplateListRequest()
	apiReq.International = common.Uint64Ptr(0)
	apiReq.Limit = common.Uint64Ptr(uint64(req.PageSize))
	apiReq.Offset = common.Uint64Ptr(uint64((req.Page - 1) * req.PageSize))

	apiResp, err := t.client.DescribeSmsTemplateList(apiReq)
	if err != nil {
		return nil, err
	}

	set := apiResp.Response.DescribeTemplateStatusSet
	templates := make([]TemplateItem, len(set))
	for i, tpl := range set {
		templates[i] = TemplateItem{
			TemplateID:      uint64ToStr(tpl.TemplateId),
			TemplateName:    strVal(tpl.TemplateName),
			TemplateContent: strVal(tpl.TemplateContent),
			Status:          tencentTmplStatus(intVal(tpl.StatusCode)),
			Reason:          strVal(tpl.ReviewReply),
			CreateTime:      unixToStr(tpl.CreateTime),
		}
	}
	return &TemplateListResponse{
		RequestID:  strVal(apiResp.Response.RequestId),
		Code:       "OK",
		Message:    "success",
		TotalCount: len(templates),
		Templates:  templates,
	}, nil
}

func (t *tencentProvider) ApplyTemplate(req *ApplyTemplateRequest) (*ApplyTemplateResponse, error) {
	apiReq := smsapi.NewAddSmsTemplateRequest()
	apiReq.TemplateName = &req.TemplateName
	apiReq.TemplateContent = &req.TemplateContent
	apiReq.Remark = &req.Remark
	apiReq.SmsType = common.Uint64Ptr(uint64(req.TemplateType + 1))
	apiReq.International = common.Uint64Ptr(0)

	apiResp, err := t.client.AddSmsTemplate(apiReq)
	if err != nil {
		return nil, err
	}
	return &ApplyTemplateResponse{
		RequestID:  strVal(apiResp.Response.RequestId),
		Code:       "OK",
		Message:    "success",
		TemplateID: strVal(apiResp.Response.AddTemplateStatus.TemplateId),
	}, nil
}

func (t *tencentProvider) SignatureList(req *SignatureListRequest) (*SignatureListResponse, error) {
	apiReq := smsapi.NewDescribeSmsSignListRequest()
	apiReq.International = common.Uint64Ptr(0)
	apiReq.Limit = common.Uint64Ptr(uint64(req.PageSize))
	apiReq.Offset = common.Uint64Ptr(uint64((req.Page - 1) * req.PageSize))

	apiResp, err := t.client.DescribeSmsSignList(apiReq)
	if err != nil {
		return nil, err
	}

	set := apiResp.Response.DescribeSignListStatusSet
	signs := make([]SignatureItem, len(set))
	for i, s := range set {
		signs[i] = SignatureItem{
			SignName:   strVal(s.SignName),
			Status:     tencentSignStatus(intVal(s.StatusCode)),
			Reason:     strVal(s.ReviewReply),
			CreateTime: unixToStr(s.CreateTime),
		}
	}
	return &SignatureListResponse{
		RequestID:  strVal(apiResp.Response.RequestId),
		Code:       "OK",
		Message:    "success",
		TotalCount: len(signs),
		Signatures: signs,
	}, nil
}

func (t *tencentProvider) ApplySignature(req *ApplySignatureRequest) (*ApplySignatureResponse, error) {
	apiReq := smsapi.NewAddSmsSignRequest()
	apiReq.SignName = &req.SignName
	apiReq.Remark = &req.Remark
	apiReq.SignType = common.Uint64Ptr(uint64(req.SignSource))
	apiReq.International = common.Uint64Ptr(0)
	apiReq.SignPurpose = common.Uint64Ptr(0)
	if req.ProofBase64 != "" {
		apiReq.ProofImage = &req.ProofBase64
	}

	apiResp, err := t.client.AddSmsSign(apiReq)
	if err != nil {
		return nil, err
	}
	return &ApplySignatureResponse{
		RequestID: strVal(apiResp.Response.RequestId),
		Code:      "OK",
		Message:   "success",
		SignName:  req.SignName,
	}, nil
}

func (t *tencentProvider) SendStatistics(req *StatisticsRequest) (*StatisticsResponse, error) {
	apiReq := smsapi.NewSendStatusStatisticsRequest()
	apiReq.BeginTime = &req.StartDate
	apiReq.EndTime = &req.EndDate
	apiReq.SmsSdkAppId = &t.appID
	apiReq.Limit = common.Uint64Ptr(0)
	apiReq.Offset = common.Uint64Ptr(0)

	apiResp, err := t.client.SendStatusStatistics(apiReq)
	if err != nil {
		return nil, err
	}

	stat := apiResp.Response.SendStatusStatistics
	return &StatisticsResponse{
		RequestID:  strVal(apiResp.Response.RequestId),
		Code:       "OK",
		Message:    "success",
		TotalSent:  int64(uint64Val(stat.RequestCount)),
		SuccessCnt: int64(uint64Val(stat.RequestSuccessCount)),
		FailCnt:    int64(uint64Val(stat.RequestCount) - uint64Val(stat.RequestSuccessCount)),
	}, nil
}

func (t *tencentProvider) QueryRecords(req *QueryRecordRequest) (*QueryRecordResponse, error) {
	apiReq := smsapi.NewPullSmsSendStatusRequest()
	apiReq.Limit = common.Uint64Ptr(uint64(req.PageSize))
	apiReq.SmsSdkAppId = &t.appID

	apiResp, err := t.client.PullSmsSendStatus(apiReq)
	if err != nil {
		return nil, err
	}

	set := apiResp.Response.PullSmsSendStatusSet
	records := make([]SendRecordItem, len(set))
	for i, r := range set {
		status := "success"
		if strVal(r.ReportStatus) == "FAIL" {
			status = "failed"
		}
		records[i] = SendRecordItem{
			PhoneNumber:  strVal(r.PhoneNumber),
			SendDate:     unixToStr(r.UserReceiveTime),
			ReceiveDate:  unixToStr(r.UserReceiveTime),
			TemplateCode: "",
			Content:      "",
			Status:       status,
			ErrCode:      strVal(r.Description),
		}
	}
	return &QueryRecordResponse{
		RequestID:  strVal(apiResp.Response.RequestId),
		Code:       "OK",
		Message:    "success",
		TotalCount: len(records),
		Records:    records,
	}, nil
}

func tencentTmplStatus(s int64) string {
	switch s {
	case 0:
		return "approved"
	case 1:
		return "pending"
	case -1:
		return "rejected"
	default:
		return ""
	}
}

func tencentSignStatus(s int64) string {
	switch s {
	case 0:
		return "approved"
	case 1:
		return "pending"
	case -1:
		return "rejected"
	default:
		return ""
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intVal(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func uint64Val(i *uint64) uint64 {
	if i == nil {
		return 0
	}
	return *i
}

func uint64ToStr(i *uint64) string {
	if i == nil {
		return ""
	}
	return strconv.FormatUint(*i, 10)
}

func unixToStr(t *uint64) string {
	if t == nil || *t == 0 {
		return ""
	}
	return time.Unix(int64(*t), 0).Format("2006-01-02 15:04:05")
}

func strPtrs(ss []string) []*string {
	ptrs := make([]*string, len(ss))
	for i := range ss {
		ptrs[i] = &ss[i]
	}
	return ptrs
}
