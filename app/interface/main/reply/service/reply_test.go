package service

import (
	"context"
	"testing"

	model "go-common/app/interface/main/reply/model/reply"
)

func TestRepressEmotion(t *testing.T) {
	got := repressEmotion("this is a test[12345]message", "[12345]")
	shouldBe := "this is a test【12345】message"
	if got != shouldBe {
		t.Fatalf("repressEmotion Error, should be: %v, got: %v", shouldBe, got)
	}
}

func TestRepressEmotions(t *testing.T) {
	got := RepressEmotions("this [8888]is a test[12345]message[4657]", []string{"[12345]", "[8888]", "[4657]"})
	shouldBe := "this 【8888】is a test【12345】message【4657】"
	if got != shouldBe {
		t.Fatalf("RepressEmotions Error, should be: %v, got: %v", shouldBe, got)
	}
}

func TestGetReportReply(t *testing.T) {
	s := &Service{}
	s.ReportReply(context.Background(), 1, 1, 1, int8(2), 3, 4, true)
}

func TestGetReplyByIDs(t *testing.T) {
	s := &Service{}
	s.GetReplyByIDs(context.Background(), int64(1), int8(1), []int64{1, 2, 3, 4})
}

func TestCheckAssist(t *testing.T) {
	s := &Service{}
	s.CheckAssist(context.Background(), int64(1), 1)
}

func TestGetRelationMap(t *testing.T) {
	s := &Service{}
	s.GetRelationMap(context.Background(), 1, []int64{}, "")
}

func TestPing(t *testing.T) {
	s := &Service{}
	s.Ping(context.Background())
}

func TestClose(t *testing.T) {
	s := &Service{}
	s.Close()
}
func TestFillRootReplies(t *testing.T) {
	s := &Service{}
	s.FillRootReplies(context.Background(),
		[]*model.Reply{},
		22,
		"123",
		true,
		&model.Subject{})
}
func TestGetBlacklist(t *testing.T) {
	s := &Service{}
	s.GetBlacklist(context.Background(), 56)
}
func TestGetFansMap(t *testing.T) {
	s := &Service{}
	s.GetFansMap(context.Background(), []int64{}, 11, "")
}
func TestGetReplyCounts(t *testing.T) {
	s := &Service{}
	s.GetReplyCounts(context.Background(), []int64{}, int8(1))
}
func TestGetAssistMap(t *testing.T) {
	s := &Service{}
	s.GetAssistMap(context.Background(), 11, "")
}
func TestUserBlockStatus(t *testing.T) {
	s := &Service{}
	s.UserBlockStatus(context.Background(), 1854)
}

func TestGetRootReplyIDs(t *testing.T) {
	s := &Service{}
	s.GetRootReplyIDs(context.Background(), 11, int8(1), int8(1), 22, 65)
}

func TestAdminReportRecover(t *testing.T) {
	s := &Service{}
	s.AdminReportRecover(context.Background(), 11, 12, 13, int8(1), int8(1), "")
}

func TestReplyContent(t *testing.T) {
	s := &Service{}
	s.ReplyContent(context.Background(), 11, 11, int8(1))
}
