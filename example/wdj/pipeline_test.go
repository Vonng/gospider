package wdj_app

import "testing"

func TestGetPipeline(t *testing.T) {
	pipe, err := GetPipeline("postgres://vonng@localhost:5432/app?sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	app, err := ParseWdjAppFromApk("com.tencent.mm")
	if err != nil {
		t.Error(err)
	}
	app.Print()

	pipe.Send(app.Item())
}
