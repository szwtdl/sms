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
    })
    resp, err := sms.SendCode(ali, "13800138000", "SMS_123456", "签名", "code", "654321")
    if err != nil {
        panic(err)
    }
    fmt.Printf("aliyun: code=%s, msg=%s\n", resp.Code, resp.Message)

    // 腾讯云
    tx, _ := sms.NewTencent(sms.TencentConfig{
        SecretID:  "your-secret-id",
        SecretKey: "your-secret-key",
        AppID:     "1400006666",
    })
    resp, err = sms.SendCode(tx, "13800138000", "1234567", "签名", "code", "654321")
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
    SignName:     "签名",
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
    SignName:     "签名",
    TemplateParams: map[string]string{"code": "654321"},
})
```

## API

### Provider 接口

```go
type Provider interface {
    Send(req *SendRequest) (*SendResponse, error)
    ProviderName() string
}
```

### SendRequest

| 字段 | 类型 | 说明 |
|------|------|------|
| PhoneNumbers | string | 手机号，多个用逗号分隔 |
| TemplateID | string | 模板 ID（云平台控制台配置） |
| TemplateParams | map[string]string | 模板参数，key 为占位符名，value 为实际值 |
| SignName | string | 短信签名（需在云平台报备） |

> **关于 TemplateParams：** 阿里云直接序列化为 JSON 对象 `{"code":"654321"}`；腾讯云按 key 字母序排列 value 后传入。单参数（验证码）场景两个平台行为一致；多参数时请确保 key 命名字母序与腾讯云模板占位符顺序一致。

### SendResponse

| 字段 | 类型 | 说明 |
|------|------|------|
| RequestID | string | 云平台请求 ID，用于排查 |
| BizID | string | 发送回执 ID |
| Code | string | 状态码，`OK` 表示成功 |
| Message | string | 状态描述 |

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

### TencentConfig

| 字段 | 说明 |
|------|------|
| SecretID | 腾讯云 SecretId |
| SecretKey | 腾讯云 SecretKey |
| AppID | 短信应用 ID（控制台 > 应用管理） |

### 便捷函数

```go
func SendCode(p Provider, phone, templateID, signName, paramKey, code string) (*SendResponse, error)
```

| 参数 | 说明 |
|------|------|
| p | Provider 实例 |
| phone | 手机号 |
| templateID | 模板 ID |
| signName | 短信签名 |
| paramKey | 模板中验证码占位符名称，如 `"code"` |
| code | 验证码 |

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
