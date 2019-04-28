package service

import (
	"flag"
	"go-common/app/job/main/up/conf"
	"go-common/app/job/main/up/dao/upcrm"
	"go-common/app/job/main/up/model/signmodel"
	"html/template"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/up-job.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	m.Run()
	os.Exit(0)
}

func TestTemplateSign(t *testing.T) {
	var data = &dueData{
		Signs: []*upcrm.SignWithName{
			{Name: "test", SignUp: signmodel.SignUp{Mid: 123, EndDate: 1540901779}},
			{Name: "tes2t", SignUp: signmodel.SignUp{Mid: 1234, EndDate: 1540902000}},
		},
	}
	tmpl, err := template.New("signTitle").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.SignTmplTitle)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}

	err = tmpl.Execute(os.Stdout, data.Signs)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}

	tmpl, err = template.New("sign").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.SignTmplContent)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	err = tmpl.Execute(os.Stdout, data.Signs)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
}

func TestTemplatePay(t *testing.T) {
	var data = &dueData{
		Pays: []*upcrm.PayWithAdmin{
			{Name: "test", SignPay: signmodel.SignPay{Mid: 123, DueDate: 1540901779, PayValue: 10000}},
			{Name: "test", SignPay: signmodel.SignPay{Mid: 123, DueDate: 1540901779, PayValue: 10000}},
			{Name: "test", SignPay: signmodel.SignPay{Mid: 123, DueDate: 1540901779, PayValue: 10000}},
		},
	}
	tmpl, err := template.New("payTitle").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.PayTmplTitle)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	err = tmpl.Execute(os.Stdout, data.Pays)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}

	tmpl, err = template.New("pay").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.PayTmplContent)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	err = tmpl.Execute(os.Stdout, data.Pays)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
}

func TestTemplateTask(t *testing.T) {
	var data = &dueData{
		Tasks: []*upcrm.TaskWithAdmin{
			{Name: "test", SignTaskHistory: signmodel.SignTaskHistory{Mid: 123, GenerateDate: 1540901779, TaskType: 2, TaskCounter: 1, TaskCondition: 10}},
			{Name: "test", SignTaskHistory: signmodel.SignTaskHistory{Mid: 123, GenerateDate: 1540901779, TaskType: 3, TaskCounter: 1, TaskCondition: 10}},
			{Name: "test", SignTaskHistory: signmodel.SignTaskHistory{Mid: 123, GenerateDate: 1540901779, TaskType: 0, TaskCounter: 1, TaskCondition: 10}},
			{Name: "test", SignTaskHistory: signmodel.SignTaskHistory{Mid: 123, GenerateDate: 1540901779, TaskType: 1, TaskCounter: 1, TaskCondition: 10}},
		}}

	tmpl, err := template.New("task").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.TaskTmplTitle)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	err = tmpl.Execute(os.Stdout, data.Tasks)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	tmpl, err = template.New("task").Funcs(funcHelper).Parse(conf.Conf.MailTemplateConf.TaskTmplContent)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
	err = tmpl.Execute(os.Stdout, data.Tasks)
	if err != nil {
		t.Errorf("err=%v", err)
		t.FailNow()
	}
}
