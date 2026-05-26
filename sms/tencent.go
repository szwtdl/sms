package sms

import (
	"sort"
	"strings"

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

func strPtrs(ss []string) []*string {
	ptrs := make([]*string, len(ss))
	for i := range ss {
		ptrs[i] = &ss[i]
	}
	return ptrs
}
