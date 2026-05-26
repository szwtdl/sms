package sms

import (
	"encoding/json"

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
