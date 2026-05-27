# sms

短信发送 Go 包，支持阿里云和腾讯云短信服务，统一接口调用。

## 安装

```bash
go get github.com/szwtdl/sms
```

## 快速开始

### 发送短信验证码（最常用）

```go
package main

import (
    "fmt"
    "github.com/szwtdl/sms/sms"
)

func main() {
    // 阿里云
    ali, _ := sms.NewAliyun(sms.AliyunConfig{
        AccessKeyID:     "your-access-key-id",
        AccessKeySecret: "your-access-key-secret",
        SignName:        "签名",
    })
    resp, err := sms.SendCode(ali, "13800138000", "SMS_123456", map[string]string{"code": "654321"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("aliyun: code=%s, msg=%s\n", resp.Code, resp.Message)

    // 腾讯云
    tx, _ := sms.NewTencent(sms.TencentConfig{
        SecretID:  "your-secret-id",
        SecretKey: "your-secret-key",
        AppID:     "1400006666",
        SignName:  "签名",
    })
    resp, err = sms.SendCode(tx, "13800138000", "1234567", map[string]string{"code": "654321"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("tencent: code=%s, msg=%s\n", resp.Code, resp.Message)
}
```

### 发送多参数模板短信

```go
ali, _ := sms.NewAliyun(sms.AliyunConfig{...})

resp, err := ali.Send(&sms.SendRequest{
    PhoneNumbers: "13800138000",
    TemplateID:   "SMS_123456",
    TemplateParams: map[string]string{
        "code":    "654321",
        "product": "XX平台",
    },
})
```

### 群发短信（多手机号逗号分隔）

```go
resp, err := ali.Send(&sms.SendRequest{
    PhoneNumbers: "13800138000,13900139000", // 阿里云逗号分隔
    TemplateID:   "SMS_123456",
    TemplateParams: map[string]string{"code": "654321"},
})
```

## API

### Provider 接口

```go
type Provider interface {
    // 短信发送
    Send(req *SendRequest) (*SendResponse, error)
    ProviderName() string

    // 模板管理
    TemplateList(req *TemplateListRequest) (*TemplateListResponse, error)
    ApplyTemplate(req *ApplyTemplateRequest) (*ApplyTemplateResponse, error)

    // 签名管理
    SignatureList(req *SignatureListRequest) (*SignatureListResponse, error)
    ApplySignature(req *ApplySignatureRequest) (*ApplySignatureResponse, error)

    // 统计与记录
    SendStatistics(req *StatisticsRequest) (*StatisticsResponse, error)
    QueryRecords(req *QueryRecordRequest) (*QueryRecordResponse, error)
}
```

---

### 短信发送

#### SendRequest

| 字段 | 类型 | 说明 |
|------|------|------|
| PhoneNumbers | string | 手机号，多个用逗号分隔 |
| TemplateID | string | 模板 ID（云平台控制台配置） |
| TemplateParams | map[string]string | 模板参数，key 为占位符名，value 为实际值 |
| SignName | string | 短信签名，传空则使用 Config 中的默认签名 |

> **关于 TemplateParams：** 阿里云直接序列化为 JSON 对象 `{"code":"654321"}`；腾讯云按 key 字母序排列 value 后传入。单参数（验证码）场景两个平台行为一致；多参数时请确保 key 命名字母序与腾讯云模板占位符顺序一致。

#### SendResponse

| 字段 | 类型 | 说明 |
|------|------|------|
| RequestID | string | 云平台请求 ID，用于排查 |
| BizID | string | 发送回执 ID |
| Code | string | 状态码，`OK` 表示成功 |
| Message | string | 状态描述 |

```go
func SendCode(p Provider, phone, templateID string, params map[string]string) (*SendResponse, error)
```

---

### 模板管理

#### 查询模板列表

```go
resp, err := p.TemplateList(&sms.TemplateListRequest{
    Page:     1,
    PageSize: 10,
})
// resp.Templates: []TemplateItem，status 为 "approved"/"pending"/"rejected"
```

#### 申请模板

```go
resp, err := p.ApplyTemplate(&sms.ApplyTemplateRequest{
    TemplateName:    "验证码模板",
    TemplateContent: "您的验证码为${code}",
    Remark:          "用于登录验证",
    TemplateType:    0, // 0=验证码, 1=通知, 2=推广
})
// resp.TemplateID: 平台分配的模板 ID（阿里云返回 SMS_xxxxxx，腾讯云返回数字 ID）
```

> 阿里云模板变量用 `${code}` 格式，腾讯云用 `{1}` 格式。

| 类型 | 字段 | 说明 |
|------|------|------|
| **请求** | TemplateName | 模板名称 |
| | TemplateContent | 模板内容 |
| | Remark | 申请说明 |
| | TemplateType | 0=验证码, 1=通知, 2=推广 |
| **响应** | TemplateID | 平台分配的模板 ID |
| | Code | `"OK"` 表示成功 |

---

### 签名管理

#### 查询签名列表

```go
resp, err := p.SignatureList(&sms.SignatureListRequest{
    Page:     1,
    PageSize: 10,
})
// resp.Signatures: []SignatureItem，status 为 "approved"/"pending"/"rejected"
```

#### 申请签名

```go
resp, err := p.ApplySignature(&sms.ApplySignatureRequest{
    SignName:    "XX科技",
    Remark:      "公司签名申请",
    SignSource:  0,           // 0=企事业单位, 1=工信部备案域名
    ProofBase64: "iVBORw0...", // 资质证明图片 base64
    ProofSuffix: "jpg",
})
```

| 类型 | 字段 | 说明 |
|------|------|------|
| **请求** | SignName | 签名名称 |
| | Remark | 申请说明 |
| | SignSource | 0=企事业单位, 1=工信部备案域名 |
| | ProofBase64 | 资质证明图片 base64 |
| | ProofSuffix | 文件后缀: jpg/png/pdf |
| **响应** | SignName | 申请的签名名称 |

---

### 发送统计与记录

#### 发送统计

```go
resp, err := p.SendStatistics(&sms.StatisticsRequest{
    StartDate: "20240501", // yyyyMMdd
    EndDate:   "20240531",
    SignName:  "", // 可选，按签名过滤
})
// resp.TotalSent / resp.SuccessCnt / resp.FailCnt
```

> 阿里云返回按日明细汇总；腾讯云返回时间段内的总量统计。

#### 查询发送记录

```go
resp, err := p.QueryRecords(&sms.QueryRecordRequest{
    PhoneNumber: "13800138000",
    SendDate:    "20240501", // yyyyMMdd
    Page:        1,
    PageSize:    10,
})
// resp.Records: []SendRecordItem，status 为 "success"/"failed"/"pending"
```

> 阿里云按 SendDate 日期查询明细；腾讯云从回调队列拉取最近的状态记录（需在控制台开通）。

---

### 构造函数

| 函数 | 说明 |
|------|------|
| `NewAliyun(AliyunConfig) (Provider, error)` | 创建阿里云短信服务 |
| `NewTencent(TencentConfig) (Provider, error)` | 创建腾讯云短信服务 |

### AliyunConfig

| 字段 | 说明 |
|------|------|
| AccessKeyID | 阿里云 AccessKey ID |
| AccessKeySecret | 阿里云 AccessKey Secret |
| SignName | 默认短信签名（需在云平台报备） |

### TencentConfig

| 字段 | 说明 |
|------|------|
| SecretID | 腾讯云 SecretId |
| SecretKey | 腾讯云 SecretKey |
| AppID | 短信应用 ID（控制台 > 应用管理） |
| SignName | 默认短信签名（需在云平台报备） |

### 预定义错误

```go
var (
    ErrMissingPhoneNumbers // 未提供手机号
    ErrMissingTemplateID   // 未提供模板 ID
    ErrMissingSignName     // 未提供短信签名
)
```

## 模板要求

- **阿里云：** 模板变量格式 `${code}`，在控制台申请模板后会得到 `SMS_xxxxxx` 格式的模板 CODE
- **腾讯云：** 模板变量格式 `{1}`，`TemplateParams` 按 key 字母序排列后传入，单参数场景无需关注顺序
