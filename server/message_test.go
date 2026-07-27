package server

import "testing"

func TestErrCode(t *testing.T) {
	err := ErrCode(20004, "个人额度不足")
	envelope := err.Envelope()
	if envelope.Code != 20004 {
		t.Fatalf("code = %d, want 20004", envelope.Code)
	}
	if envelope.Message != "个人额度不足" {
		t.Fatalf("message = %q, want %q", envelope.Message, "个人额度不足")
	}
}

func TestErrMsgUsesDefaultCode(t *testing.T) {
	if code := ErrMsg("失败").Envelope().Code; code != -1 {
		t.Fatalf("code = %d, want -1", code)
	}
}
