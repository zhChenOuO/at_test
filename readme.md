# Amazing Talker 題目

## 說明
1. 內有兩個服務 一個是 api server, 一個是排程
2. 排程是用於 寄送 email & sms 驗證
3. Config 主要由 deploy/config/app.yaml 
4. app.yaml 需填上
5. deploy/config/credentials.json 需填上
6. token.json 需填上 
7. makefile 有完整的啟動指令

## API

### /v1/user/register

```json
{
	"register_type":1,
	"email":"test",
	"name":"test",
	"phone":"12345678",
	"phone_area_code":"886",
	"password":"12345678",
	"verify_password":"12345678"
}
```
#### Request:
|欄位|說明|型態|備註|
|-|-|-|-|
|register_type|註冊模式 1:email,2:phone|int|
|email|電子信箱|string|註冊模式1 必填|
|name|姓名|string|
|phone|電話號碼|string|註冊模式1 必填|
|phone_area_code|電話區碼|string|註冊模式1 必填|
|password|密碼|string|
|verify_password|確認密碼|string|

#### Response:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6InRlc3QiLCJlbWFpbCI6IndvcHMwMzE2QGdtYWlsLmNvbSIsInBob25lIjoiIiwiZXhwaXJlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiZXhwIjoxNjM5NjAzMTUwfQ.UKGE3T9WssN3hcy3raJVRNHSXLyPtBzj-pIoul3RD6s"
}
```