package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/dao"
	"go-common/app/service/main/spy/model"
	"go-common/library/ecode"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestServiceNewJudgementInfo(t *testing.T) {
	convey.Convey("NewJudgementInfo", t, func(ctx convey.C) {
		var (
			s   = &Service{}
			c   = context.Background()
			mid = int64(0)
			ip  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewJudgementInfo(c, s, mid, ip)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetAuditInfo(t *testing.T) {
	var (
		mid       = int64(0)
		ip        = ""
		ji        = NewJudgementInfo(context.TODO(), s, mid, ip)
		auditInfo = &model.AuditInfo{
			Mid:      mid,
			BindTel:  true,
			BindMail: false,
		}
	)

	convey.Convey("getAuditInfo", t, func(ctx convey.C) {
		var testCases = []struct {
			jiAuditInfo     *model.AuditInfo
			daoAuditInfo    *model.AuditInfo
			expectAuditInfo *model.AuditInfo
		}{
			{
				jiAuditInfo:     auditInfo,
				daoAuditInfo:    nil,
				expectAuditInfo: auditInfo,
			},
			{
				jiAuditInfo:     nil,
				daoAuditInfo:    auditInfo,
				expectAuditInfo: auditInfo,
			},
			{
				jiAuditInfo:     nil,
				daoAuditInfo:    nil,
				expectAuditInfo: nil,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d", idx), func(ctx convey.C) {
				monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "AuditInfo", func(_ *dao.Dao, _ context.Context, _ int64, _ string) (*model.AuditInfo, error) {
					if testCase.daoAuditInfo == nil {
						return nil, errors.New("s.dao.AuditInfo error")
					}
					return testCase.daoAuditInfo, nil
				})
				ji.auditInfo = testCase.jiAuditInfo
				ai, err := ji.getAuditInfo()
				if testCase.expectAuditInfo == nil {
					ctx.So(err, convey.ShouldBeError)
					ctx.So(ai, convey.ShouldBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(ai, convey.ShouldEqual, testCase.expectAuditInfo)
				}
			})
		}

	})
}

func TestServicegetProfileInfo(t *testing.T) {
	var (
		mid         = int64(0)
		ip          = ""
		ji          = NewJudgementInfo(context.TODO(), s, mid, ip)
		profileInfo = &model.ProfileInfo{
			Mid:            mid,
			Identification: 0,
		}
	)

	convey.Convey("getProfileInfo", t, func(ctx convey.C) {
		var testCases = []struct {
			jiProfileInfo     *model.ProfileInfo
			daoProfileInfo    *model.ProfileInfo
			expectProfileInfo *model.ProfileInfo
		}{
			{
				jiProfileInfo:     profileInfo,
				daoProfileInfo:    nil,
				expectProfileInfo: profileInfo,
			},
			{
				jiProfileInfo:     nil,
				daoProfileInfo:    profileInfo,
				expectProfileInfo: profileInfo,
			},
			{
				jiProfileInfo:     nil,
				daoProfileInfo:    nil,
				expectProfileInfo: nil,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d", idx), func(ctx convey.C) {
				monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "ProfileInfo", func(_ *dao.Dao, _ context.Context, _ int64, _ string) (*model.ProfileInfo, error) {
					if testCase.daoProfileInfo == nil {
						return nil, errors.New("s.dao.ProfileInfo error")
					}
					return testCase.daoProfileInfo, nil
				})
				ji.profileInfo = testCase.jiProfileInfo
				pi, err := ji.getProfileInfo()
				if testCase.expectProfileInfo == nil {
					ctx.So(err, convey.ShouldBeError)
					ctx.So(pi, convey.ShouldBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(pi, convey.ShouldEqual, testCase.expectProfileInfo)
				}
			})
		}

	})
}

func TestServicegetTelRiskInfo(t *testing.T) {
	var (
		mid                 = int64(0)
		ip                  = ""
		ji                  = NewJudgementInfo(context.TODO(), s, mid, ip)
		telRiskLevel        = int8(1)
		restoreEventHistory = &model.UserEventHistory{
			Mid:     mid,
			EventID: 1024,
		}
		telRiskInfo = &model.TelRiskInfo{
			TelRiskLevel:   telRiskLevel,
			RestoreHistory: restoreEventHistory,
		}
	)

	convey.Convey("getTelRiskInfo", t, func(ctx convey.C) {
		var testCases = []struct {
			jiTelRiskInfo          *model.TelRiskInfo
			daoTelRiskLevel        int8
			daoRestoreEventHistory *model.UserEventHistory
			expectTelRiskInfo      *model.TelRiskInfo
		}{
			{
				jiTelRiskInfo:          telRiskInfo,
				daoTelRiskLevel:        int8(0),
				daoRestoreEventHistory: nil,
				expectTelRiskInfo:      telRiskInfo,
			},
			{
				jiTelRiskInfo:          nil,
				daoTelRiskLevel:        telRiskLevel,
				daoRestoreEventHistory: restoreEventHistory,
				expectTelRiskInfo:      telRiskInfo,
			},
			{
				jiTelRiskInfo:          nil,
				daoTelRiskLevel:        int8(0),
				daoRestoreEventHistory: nil,
				expectTelRiskInfo:      nil,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d", idx), func(ctx convey.C) {
				monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "TelRiskLevel", func(_ *dao.Dao, _ context.Context, _ int64) (int8, error) {
					return testCase.daoTelRiskLevel, nil
				})
				monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "Event", func(_ *dao.Dao, _ context.Context, _ string) (*model.Event, error) {
					return &model.Event{
						ID: mid,
					}, nil
				})

				monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "EventHistoryByMidAndEvent", func(_ *dao.Dao, _ context.Context, _ int64, _ int64) (*model.UserEventHistory, error) {
					if testCase.daoRestoreEventHistory == nil {
						return nil, errors.New("s.dao.EventHistoryByMidAndEvent error")
					}
					return testCase.daoRestoreEventHistory, nil
				})

				ji.telRiskInfo = testCase.jiTelRiskInfo
				tri, err := ji.getTelRiskInfo()
				if testCase.expectTelRiskInfo == nil {
					ctx.So(err, convey.ShouldBeError)
					ctx.So(tri, convey.ShouldBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(tri, convey.ShouldResemble, testCase.expectTelRiskInfo)
				}
			})

		}

	})
}

func TestServicegetRulesMap(t *testing.T) {
	convey.Convey("getRulesMap", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := s.getRulesMap()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetJudgementInfo(t *testing.T) {
	convey.Convey("getJudgementInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ip  = ""
		)
		convey.Convey("getJudgementInfo", func(ctx convey.C) {
			ji, err := s.getJudgementInfo(c, mid, ip)
			ctx.So(ji, convey.ShouldNotBeNil)
			ctx.So(ji.auditInfo, convey.ShouldBeNil)
			ctx.So(ji.profileInfo, convey.ShouldBeNil)
			ctx.So(ji.telRiskInfo, convey.ShouldBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})

	})
}

func TestServicegetRuleFunc(t *testing.T) {
	convey.Convey("getRuleFunc", t, func(ctx convey.C) {
		var testCases = []struct {
			rule         Rule
			expectedResp RuleFunc
		}{
			{
				rule:         _telIsBound,
				expectedResp: telIsBound,
			},
			{
				rule:         "Anything",
				expectedResp: nil,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d..,", idx), func(ctx convey.C) {
				ruleFunc, err := s.getRuleFunc(testCase.rule)
				if testCase.expectedResp == nil {
					ctx.So(err, convey.ShouldBeError, ecode.SpyRuleNotExist)
					ctx.So(ruleFunc, convey.ShouldBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(ruleFunc, convey.ShouldEqual, testCase.expectedResp)
				}
			})
		}

	})
}

func TestServiceiterRules(t *testing.T) {
	var (
		rules = Rules{
			_telIsBound, _mailIsBound,
		}
	)

	convey.Convey("iterRules", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectRes     IterState
		}{
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel:  true,
						BindMail: true,
					},
				},
				expectRes: IterSuccess,
			},
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel:  true,
						BindMail: false,
					},
				},
				expectRes: IterFail,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				state, err := s.iterRules(testCase.judgementInfo, rules)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldEqual, testCase.expectRes)
			})
		}
	})
}

func TestServicegetBaseScoreFactor(t *testing.T) {
	convey.Convey("getBaseScoreFactor", t, func(ctx convey.C) {
		var (
			testCases = []struct {
				judgementInfo *JudgementInfo
				expectRest    FactoryMeta
			}{
				{
					judgementInfo: &JudgementInfo{
						auditInfo: &model.AuditInfo{
							BindTel:  false,
							BindMail: true,
						},
						telRiskInfo: &model.TelRiskInfo{
							RestoreHistory: &model.UserEventHistory{
								Mid: int64(0),
							},
						},
					},
					expectRest: FactoryMeta{
						serviceName: _defaultService,
						eventName:   conf.Conf.Property.Event.BindMailAndIdenUnknown,
						riskLevel:   _defaultRiskLevel,
					},
				},
				{
					judgementInfo: &JudgementInfo{
						auditInfo: &model.AuditInfo{
							BindTel:  false,
							BindMail: false,
						},
						profileInfo: &model.ProfileInfo{
							Identification: 1,
						},
						telRiskInfo: &model.TelRiskInfo{
							RestoreHistory: &model.UserEventHistory{
								Mid: int64(0),
							},
						},
					},
					expectRest: FactoryMeta{},
				},
			}
		)
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				meta, err := s.getBaseScoreFactor(testCase.judgementInfo)
				if testCase.expectRest == (FactoryMeta{}) {
					ctx.So(err, convey.ShouldBeError, ecode.SpyRulesNotMatch)
					ctx.So(meta, convey.ShouldResemble, FactoryMeta{})
				} else {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(meta, convey.ShouldResemble, testCase.expectRest)
				}
			})
		}
	})
}

func TestServicetelIsBound(t *testing.T) {
	convey.Convey("telIsBound", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel: true,
					},
				},
				expectResp: true,
			}, {
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel: false,
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsBound(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})
}

func TestServicetelIsNotBound(t *testing.T) {
	convey.Convey("telIsNotBound", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel: true,
					},
				},
				expectResp: false,
			}, {
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindTel: false,
					},
				},
				expectResp: true,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsNotBound(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})
}

func TestServicetelIsLowRisk(t *testing.T) {
	convey.Convey("telIsLowRisk", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:   dao.TelRiskLevelLow,
						RestoreHistory: nil,
					},
				},
				expectResp: true,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel: dao.TelRiskLevelHigh,
						RestoreHistory: &model.UserEventHistory{
							Mid: int64(0),
						},
					},
				},
				expectResp: true,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:   dao.TelRiskLevelHigh,
						RestoreHistory: nil,
					},
				},
				expectResp: false,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:    dao.TelRiskLevelHigh,
						RestoreHistory:  nil,
						UnicomGiftState: 1,
					},
				},
				expectResp: true,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsLowRisk(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})
}

func TestServicetelIsMediumRisk(t *testing.T) {
	convey.Convey("telIsMediumRisk", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:   dao.TelRiskLevelMedium,
						RestoreHistory: nil,
					},
				},
				expectResp: true,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel: dao.TelRiskLevelMedium,
						RestoreHistory: &model.UserEventHistory{
							Mid: int64(0),
						},
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsMediumRisk(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})
}

func TestServicetelIsHighRisk(t *testing.T) {
	convey.Convey("telIsHighRisk", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:   dao.TelRiskLevelHigh,
						RestoreHistory: nil,
					},
				},
				expectResp: true,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:    dao.TelRiskLevelHigh,
						RestoreHistory:  nil,
						UnicomGiftState: 1,
					},
				},
				expectResp: false,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel: dao.TelRiskLevelMedium,
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsHighRisk(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})

}

func TestServicetelIsUnknownRisk(t *testing.T) {

	convey.Convey("telIsUnknownRisk", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel:   dao.TelRiskLevelUnknown,
						RestoreHistory: nil,
					},
				},
				expectResp: true,
			},
			{
				judgementInfo: &JudgementInfo{
					telRiskInfo: &model.TelRiskInfo{
						TelRiskLevel: dao.TelRiskLevelMedium,
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := telIsUnknownRisk(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})

}

func TestServicemailIsBound(t *testing.T) {
	convey.Convey("mailIsBound", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindMail: true,
					},
				},
				expectResp: true,
			}, {
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindMail: false,
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := mailIsBound(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})

}

func TestServicemailIsNotBound(t *testing.T) {
	convey.Convey("mailIsNotBound", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindMail: true,
					},
				},
				expectResp: false,
			}, {
				judgementInfo: &JudgementInfo{
					auditInfo: &model.AuditInfo{
						BindMail: false,
					},
				},
				expectResp: true,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := mailIsNotBound(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})

}

func TestServiceidenIsNotAuthed(t *testing.T) {
	convey.Convey("idenIsNotAuthed", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					profileInfo: &model.ProfileInfo{
						Identification: 0,
					},
				},
				expectResp: true,
			}, {
				judgementInfo: &JudgementInfo{
					profileInfo: &model.ProfileInfo{
						Identification: 1,
					},
				},
				expectResp: false,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := idenIsNotAuthed(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})

}

func TestServiceidenIsAuthed(t *testing.T) {
	convey.Convey("idenIsAuthed", t, func(ctx convey.C) {
		var testCases = []struct {
			judgementInfo *JudgementInfo
			expectResp    bool
		}{
			{
				judgementInfo: &JudgementInfo{
					profileInfo: &model.ProfileInfo{
						Identification: 0,
					},
				},
				expectResp: false,
			}, {
				judgementInfo: &JudgementInfo{
					profileInfo: &model.ProfileInfo{
						Identification: 1,
					},
				},
				expectResp: true,
			},
		}
		for idx, testCase := range testCases {
			ctx.Convey(fmt.Sprintf("Iterating Case%d...", idx), func(ctx convey.C) {
				result, err := idenIsAuthed(testCase.judgementInfo)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, testCase.expectResp)
			})
		}

	})
}
