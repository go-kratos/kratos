package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/upcredit/conf:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business_score.go",
        "calculate_service.go",
        "content_ext.go",
        "creditscore.go",
        "databus.go",
        "logcredit.go",
        "service.go",
    ],
    importmap = "go-common/app/service/main/upcredit/service",
    importpath = "go-common/app/service/main/upcredit/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/upcredit/common/election:go_default_library",
        "//app/service/main/upcredit/common/fsm:go_default_library",
        "//app/service/main/upcredit/conf:go_default_library",
        "//app/service/main/upcredit/dao/upcrmdao:go_default_library",
        "//app/service/main/upcredit/mathutil:go_default_library",
        "//app/service/main/upcredit/model/calculator:go_default_library",
        "//app/service/main/upcredit/model/canal:go_default_library",
        "//app/service/main/upcredit/model/upcrmmodel:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
