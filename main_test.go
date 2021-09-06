package main

import (
	"testing" // 单元测试文件需要导入testing包
)


func TestSendMail(t *testing.T) { // 函数名必须以Test开头后面加测试的函数名（函数名首字母大写），变量和变量类型是固定的
	mailTo := []string{
		"3238813014@qq.com",
	}
	subject := "Hello by golang gomail from exmail.qq.com"
	body := "Hello,by gomail sent"
	err := SendMail(mailTo, subject, body)
	if err != nil {
		t.Fatalf("SendMail function test error: %v", err)
		return
	}
	t.Log("SendMail function test : successful.")
}

func TestSqlClient_NewSubscribe(t *testing.T) {
	// Init sql client
	sqlClient := &SqlClient{}
	err := sqlClient.Init()
	if err != nil{
		t.Fatalf("SqlClient.Init function test error: %v", err)
	}
	subscribe := &SubscribeModel{
		ProposalId: "test3",
		Email: "3238813014@qq.com",
	}
	err = sqlClient.NewSubscribe(*subscribe)
	if err != nil{
		t.Fatalf("NewSubscribe function test error: %v", err)
	}
	t.Log("NewSubscribe function test : successful.")
}

func TestSqlClient_GetSubscribeEmail(t *testing.T) {
	sqlClient := &SqlClient{}
	err := sqlClient.Init()
	if err != nil{
		t.Fatalf("SqlClient.Init function test error: %v", err)
	}
	data, _ := sqlClient.GetSubscribeEmail("QmYbrvaEQhNvVfTZqNydiewJkfycvZwUDqFLc8xSPPeGVD")
	t.Logf("Test success, emails are:%v", data)
}

func TestRenderMailTmpl(t *testing.T) {
	var p PropRender
	pp, err := RenderMailTmpl(p)
	t.Log(err)
	t.Log(pp)
}