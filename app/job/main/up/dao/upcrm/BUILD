package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = ["//app/job/main/up/conf:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "sign_up.go",
        "task_info.go",
        "up_base_info.go",
    ],
    importpath = "go-common/app/job/main/up/dao/upcrm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/util/now:go_default_library",
        "//app/job/main/up/conf:go_default_library",
        "//app/job/main/up/model/signmodel:go_default_library",
        "//app/job/main/up/model/upcrmmodel:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
