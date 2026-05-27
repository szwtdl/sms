package sms

import (
	"testing"
)

// mockProvider 用于测试的 mock 实现。
type mockProvider struct {
	name     string
	sendFunc func(req *SendRequest) (*SendResponse, error)
}

func (m *mockProvider) ProviderName() string { return m.name }

func (m *mockProvider) Send(req *SendRequest) (*SendResponse, error) {
	return m.sendFunc(req)
}

func (m *mockProvider) TemplateList(req *TemplateListRequest) (*TemplateListResponse, error) {
	return &TemplateListResponse{Code: "OK", TotalCount: 0}, nil
}

func (m *mockProvider) ApplyTemplate(req *ApplyTemplateRequest) (*ApplyTemplateResponse, error) {
	return &ApplyTemplateResponse{Code: "OK", TemplateID: "TPL_NEW"}, nil
}

func (m *mockProvider) ModifyTemplate(req *ModifyTemplateRequest) (*ModifyTemplateResponse, error) {
	return &ModifyTemplateResponse{Code: "OK", TemplateID: req.TemplateID}, nil
}

func (m *mockProvider) DeleteTemplate(req *DeleteTemplateRequest) (*DeleteTemplateResponse, error) {
	return &DeleteTemplateResponse{Code: "OK"}, nil
}

func (m *mockProvider) SignatureList(req *SignatureListRequest) (*SignatureListResponse, error) {
	return &SignatureListResponse{Code: "OK", TotalCount: 0}, nil
}

func (m *mockProvider) ApplySignature(req *ApplySignatureRequest) (*ApplySignatureResponse, error) {
	return &ApplySignatureResponse{Code: "OK", SignName: req.SignName}, nil
}

func (m *mockProvider) ModifySignature(req *ModifySignatureRequest) (*ModifySignatureResponse, error) {
	return &ModifySignatureResponse{Code: "OK", SignName: req.SignName}, nil
}

func (m *mockProvider) DeleteSignature(req *DeleteSignatureRequest) (*DeleteSignatureResponse, error) {
	return &DeleteSignatureResponse{Code: "OK"}, nil
}

func (m *mockProvider) SendStatistics(req *StatisticsRequest) (*StatisticsResponse, error) {
	return &StatisticsResponse{Code: "OK", TotalSent: 100, SuccessCnt: 95, FailCnt: 5}, nil
}

func (m *mockProvider) QueryRecords(req *QueryRecordRequest) (*QueryRecordResponse, error) {
	return &QueryRecordResponse{Code: "OK", TotalCount: 0}, nil
}

func TestSendCode(t *testing.T) {
	mock := &mockProvider{
		name: "mock",
		sendFunc: func(req *SendRequest) (*SendResponse, error) {
			if req.PhoneNumbers != "13800138000" {
				t.Errorf("unexpected phone: %s", req.PhoneNumbers)
			}
			if req.TemplateID != "TPL_001" {
				t.Errorf("unexpected template id: %s", req.TemplateID)
			}
			if req.SignName != "" {
				t.Errorf("unexpected sign name: %s", req.SignName)
			}
			if len(req.TemplateParams) != 1 || req.TemplateParams["code"] != "123456" {
				t.Errorf("unexpected template params: %v", req.TemplateParams)
			}
			return &SendResponse{Code: "OK", Message: "success"}, nil
		},
	}

	resp, err := SendCode(mock, "13800138000", "TPL_001", map[string]string{"code": "123456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *SendRequest
		wantErr error
	}{
		{"missing phone", &SendRequest{TemplateID: "1", SignName: "S"}, ErrMissingPhoneNumbers},
		{"missing template", &SendRequest{PhoneNumbers: "1", SignName: "S"}, ErrMissingTemplateID},
		{"missing sign", &SendRequest{PhoneNumbers: "1", TemplateID: "1"}, ErrMissingSignName},
		{"all ok", &SendRequest{PhoneNumbers: "1", TemplateID: "1", SignName: "S"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.req)
			if err != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestSplitPhones(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"13800138000", []string{"13800138000"}},
		{"13800138000,13800138001", []string{"13800138000", "13800138001"}},
		{" 13800138000 , 13800138001 ", []string{"13800138000", "13800138001"}},
		{"13800138000, ", []string{"13800138000"}},
		{" , ", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitPhones(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("splitPhones(%q) = %v (len=%d), want %v (len=%d)", tt.input, got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitPhones(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestStrPtrs(t *testing.T) {
	got := strPtrs([]string{"a", "b"})
	if len(got) != 2 || *got[0] != "a" || *got[1] != "b" {
		t.Errorf("unexpected strPtrs result: %v", got)
	}
}

func TestMapValuesSorted(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]string
		want []string
	}{
		{"empty", map[string]string{}, nil},
		{"single", map[string]string{"code": "123456"}, []string{"123456"}},
		{"sorted", map[string]string{"b": "2", "a": "1"}, []string{"1", "2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapValuesSorted(tt.m)
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if *got[i] != tt.want[i] {
					t.Errorf("[%d] = %q, want %q", i, *got[i], tt.want[i])
				}
			}
		})
	}
}

func TestProviderName(t *testing.T) {
	mock := &mockProvider{name: "test-provider"}
	if mock.ProviderName() != "test-provider" {
		t.Errorf("unexpected provider name: %s", mock.ProviderName())
	}
}

func TestMapValuesSorted_Deterministic(t *testing.T) {
	m := map[string]string{"z": "last", "a": "first", "m": "middle"}
	got := mapValuesSorted(m)
	want := []string{"first", "middle", "last"}
	for i := range got {
		if *got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, *got[i], want[i])
		}
	}
}

func TestSendCode_MultiParams(t *testing.T) {
	p := &mockProvider{
		name: "mock",
		sendFunc: func(req *SendRequest) (*SendResponse, error) {
			if req.TemplateParams["code"] != "123456" {
				t.Errorf("unexpected code: %s", req.TemplateParams["code"])
			}
			if req.TemplateParams["product"] != "test" {
				t.Errorf("unexpected product: %s", req.TemplateParams["product"])
			}
			return &SendResponse{Code: "OK"}, nil
		},
	}

	req := &SendRequest{
		PhoneNumbers: "13800138000",
		TemplateID:   "TPL_002",
		SignName:     "测试",
		TemplateParams: map[string]string{
			"code":    "123456",
			"product": "test",
		},
	}
	resp, err := p.Send(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestMockTemplateList(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.TemplateList(&TemplateListRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestMockApplyTemplate(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.ApplyTemplate(&ApplyTemplateRequest{
		TemplateName:    "验证码模板",
		TemplateContent: "您的验证码为{1}",
		TemplateType:    0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" || resp.TemplateID != "TPL_NEW" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestMockSignatureList(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.SignatureList(&SignatureListRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestMockApplySignature(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.ApplySignature(&ApplySignatureRequest{
		SignName:    "测试签名",
		Remark:      "测试申请",
		SignSource:  0,
		ProofBase64: "base64content",
		ProofSuffix: "jpg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" || resp.SignName != "测试签名" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestMockSendStatistics(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.SendStatistics(&StatisticsRequest{
		StartDate: "20240501",
		EndDate:   "20240531",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalSent != 100 || resp.SuccessCnt != 95 || resp.FailCnt != 5 {
		t.Errorf("unexpected stats: total=%d success=%d fail=%d", resp.TotalSent, resp.SuccessCnt, resp.FailCnt)
	}
}

func TestMockQueryRecords(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.QueryRecords(&QueryRecordRequest{
		PhoneNumber: "13800138000",
		SendDate:    "20240501",
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestMockModifyTemplate(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.ModifyTemplate(&ModifyTemplateRequest{
		TemplateID:      "SMS_001",
		TemplateName:    "新验证码模板",
		TemplateContent: "您的验证码为${code}",
		TemplateType:    0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" || resp.TemplateID != "SMS_001" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestMockModifySignature(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.ModifySignature(&ModifySignatureRequest{
		SignID:      12345,
		SignName:    "新签名",
		Remark:      "修改签名",
		SignSource:  0,
		ProofBase64: "newbase64",
		ProofSuffix: "jpg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" || resp.SignName != "新签名" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestMockDeleteTemplate(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.DeleteTemplate(&DeleteTemplateRequest{
		TemplateID: "SMS_001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}

func TestMockDeleteSignature(t *testing.T) {
	m := &mockProvider{name: "mock"}
	resp, err := m.DeleteSignature(&DeleteSignatureRequest{
		SignID:   12345,
		SignName: "测试签名",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "OK" {
		t.Errorf("expected OK, got %s", resp.Code)
	}
}
