package sms

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// AliyunConfig 阿里云短信配置。
type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string // 默认短信签名，Send 时 SignName 为空则使用此值
}

type aliyunProvider struct {
	client   *dysmsapi.Client
	signName string
}

// NewAliyun 创建阿里云短信服务。
func NewAliyun(cfg AliyunConfig) (Provider, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &aliyunProvider{client: client, signName: cfg.SignName}, nil
}

func (a *aliyunProvider) ProviderName() string { return "aliyun" }

func (a *aliyunProvider) Send(req *SendRequest) (*SendResponse, error) {
	if req.SignName == "" {
		req.SignName = a.signName
	}
	if err := validate(req); err != nil {
		return nil, err
	}

	tplJSON, _ := json.Marshal(req.TemplateParams)

	apiReq := dysmsapi.CreateSendSmsRequest()
	apiReq.PhoneNumbers = req.PhoneNumbers
	apiReq.SignName = req.SignName
	apiReq.TemplateCode = req.TemplateID
	apiReq.TemplateParam = string(tplJSON)

	apiResp, err := a.client.SendSms(apiReq)
	if err != nil {
		return nil, err
	}

	resp := &SendResponse{
		RequestID: apiResp.RequestId,
		BizID:     apiResp.BizId,
		Code:      apiResp.Code,
		Message:   apiResp.Message,
	}
	return resp, nil
}

func (a *aliyunProvider) TemplateList(req *TemplateListRequest) (*TemplateListResponse, error) {
	apiReq := dysmsapi.CreateQuerySmsTemplateListRequest()
	apiReq.PageIndex = requests.NewInteger(req.Page)
	apiReq.PageSize = requests.NewInteger(req.PageSize)

	apiResp, err := a.client.QuerySmsTemplateList(apiReq)
	if err != nil {
		return nil, err
	}

	templates := make([]TemplateItem, len(apiResp.SmsTemplateList))
	for i, t := range apiResp.SmsTemplateList {
		templates[i] = TemplateItem{
			TemplateID:      t.TemplateCode,
			TemplateName:    t.TemplateName,
			TemplateContent: t.TemplateContent,
			Status:          auditStatusToStr(t.AuditStatus),
			Reason:          t.Reason.RejectInfo,
			CreateTime:      t.CreateDate,
		}
	}
	return &TemplateListResponse{
		RequestID:  apiResp.RequestId,
		Code:       apiResp.Code,
		Message:    apiResp.Message,
		TotalCount: int(apiResp.TotalCount),
		Templates:  templates,
	}, nil
}

func (a *aliyunProvider) ApplyTemplate(req *ApplyTemplateRequest) (*ApplyTemplateResponse, error) {
	apiReq := dysmsapi.CreateAddSmsTemplateRequest()
	apiReq.TemplateName = req.TemplateName
	apiReq.TemplateContent = req.TemplateContent
	apiReq.Remark = req.Remark
	apiReq.TemplateType = requests.NewInteger(req.TemplateType)

	apiResp, err := a.client.AddSmsTemplate(apiReq)
	if err != nil {
		return nil, err
	}
	return &ApplyTemplateResponse{
		RequestID:  apiResp.RequestId,
		Code:       apiResp.Code,
		Message:    apiResp.Message,
		TemplateID: apiResp.TemplateCode,
	}, nil
}

func (a *aliyunProvider) SignatureList(req *SignatureListRequest) (*SignatureListResponse, error) {
	apiReq := dysmsapi.CreateQuerySmsSignListRequest()
	apiReq.PageIndex = requests.NewInteger(req.Page)
	apiReq.PageSize = requests.NewInteger(req.PageSize)

	apiResp, err := a.client.QuerySmsSignList(apiReq)
	if err != nil {
		return nil, err
	}

	signs := make([]SignatureItem, len(apiResp.SmsSignList))
	for i, s := range apiResp.SmsSignList {
		signs[i] = SignatureItem{
			SignName:   s.SignName,
			Status:     auditStatusToStr(s.AuditStatus),
			Reason:     s.Reason.RejectInfo,
			CreateTime: s.CreateDate,
		}
	}
	return &SignatureListResponse{
		RequestID:  apiResp.RequestId,
		Code:       apiResp.Code,
		Message:    apiResp.Message,
		TotalCount: int(apiResp.TotalCount),
		Signatures: signs,
	}, nil
}

func (a *aliyunProvider) ApplySignature(req *ApplySignatureRequest) (*ApplySignatureResponse, error) {
	apiReq := dysmsapi.CreateAddSmsSignRequest()
	apiReq.SignName = req.SignName
	apiReq.Remark = req.Remark
	apiReq.SignSource = requests.NewInteger(req.SignSource)
	if req.ProofBase64 != "" {
		apiReq.SignFileList = &[]dysmsapi.AddSmsSignSignFileList{
			{FileContents: req.ProofBase64, FileSuffix: req.ProofSuffix},
		}
	}

	apiResp, err := a.client.AddSmsSign(apiReq)
	if err != nil {
		return nil, err
	}
	return &ApplySignatureResponse{
		RequestID: apiResp.RequestId,
		Code:      apiResp.Code,
		Message:   apiResp.Message,
		SignName:  apiResp.SignName,
	}, nil
}

func (a *aliyunProvider) SendStatistics(req *StatisticsRequest) (*StatisticsResponse, error) {
	apiReq := dysmsapi.CreateQuerySendStatisticsRequest()
	apiReq.IsGlobe = requests.NewInteger(req.Type)
	apiReq.StartDate = req.StartDate
	apiReq.EndDate = req.EndDate
	apiReq.SignName = req.SignName
	apiReq.PageIndex = requests.NewInteger(req.Page)
	apiReq.PageSize = requests.NewInteger(req.PageSize)

	apiResp, err := a.client.QuerySendStatistics(apiReq)
	if err != nil {
		return nil, err
	}

	var totalSent, success, fail int64
	for _, t := range apiResp.Data.TargetList {
		totalSent += t.TotalCount
		success += t.RespondedSuccessCount
		fail += t.RespondedFailCount
	}
	return &StatisticsResponse{
		RequestID:  apiResp.RequestId,
		Code:       apiResp.Code,
		Message:    apiResp.Message,
		TotalSent:  totalSent,
		SuccessCnt: success,
		FailCnt:    fail,
	}, nil
}

func (a *aliyunProvider) QueryRecords(req *QueryRecordRequest) (*QueryRecordResponse, error) {
	apiReq := dysmsapi.CreateQuerySendDetailsRequest()
	apiReq.PhoneNumber = req.PhoneNumber
	apiReq.SendDate = req.SendDate
	apiReq.PageSize = requests.NewInteger(req.PageSize)
	apiReq.CurrentPage = requests.NewInteger(req.Page)

	apiResp, err := a.client.QuerySendDetails(apiReq)
	if err != nil {
		return nil, err
	}

	records := make([]SendRecordItem, len(apiResp.SmsSendDetailDTOs.SmsSendDetailDTO))
	for i, d := range apiResp.SmsSendDetailDTOs.SmsSendDetailDTO {
		records[i] = SendRecordItem{
			PhoneNumber:  d.PhoneNum,
			SendDate:     d.SendDate,
			ReceiveDate:  d.ReceiveDate,
			TemplateCode: d.TemplateCode,
			Content:      d.Content,
			Status:       sendStatusToStr(d.SendStatus),
			ErrCode:      d.ErrCode,
		}
	}
	return &QueryRecordResponse{
		RequestID:  apiResp.RequestId,
		Code:       apiResp.Code,
		Message:    apiResp.Message,
		TotalCount: atoi(apiResp.TotalCount),
		Records:    records,
	}, nil
}

func auditStatusToStr(status string) string {
	switch status {
	case "AUDIT_STATE_PASS", "PASS", "approved":
		return "approved"
	case "AUDIT_STATE_NOT_PASS", "NOT_PASS", "rejected":
		return "rejected"
	case "AUDIT_STATE_AUDITING", "AUDITING", "pending":
		return "pending"
	default:
		return status
	}
}

func sendStatusToStr(status int64) string {
	switch status {
	case 1:
		return "pending"
	case 2:
		return "failed"
	case 3:
		return "success"
	default:
		return ""
	}
}

func atoi(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func validate(req *SendRequest) error {
	if req.PhoneNumbers == "" {
		return ErrMissingPhoneNumbers
	}
	if req.TemplateID == "" {
		return ErrMissingTemplateID
	}
	if req.SignName == "" {
		return ErrMissingSignName
	}
	return nil
}
