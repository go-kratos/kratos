package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["arcfsm_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/upcredit/conf:go_default_library",
        "//app/service/main/upcredit/model/upcrmmodel:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "creditscore.go",
        "score_rules.go",
    ],
    importmap = "go-common/app/service/main/upcredit/model/calculator",
    importpath = "go-common/app/service/main/upcredit/model/calculator",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/upcredit/common/fsm:go_default_library",
        "//app/service/main/upcredit/conf:go_default_library",
        "//app/service/main/upcredit/mathutil:go_default_library",
        "//app/service/main/upcredit/model/upcrmmodel:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
