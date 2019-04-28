package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "associate_outer_test.go",
        "coupon_test.go",
        "manager_test.go",
        "panel_test.go",
        "vip_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/account/conf:go_default_library",
        "//app/interface/main/account/model:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "associate.go",
        "associate_limit.go",
        "associate_outer.go",
        "coupon.go",
        "manager.go",
        "panel.go",
        "privilege.go",
        "resource.go",
        "tips.go",
        "vip.go",
        "welfare.go",
    ],
    importpath = "go-common/app/interface/main/account/service/vip",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/account/conf:go_default_library",
        "//app/interface/main/account/dao/vip:go_default_library",
        "//app/interface/main/account/model:go_default_library",
        "//app/service/main/coupon/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/rpc/client:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//app/service/main/resource/rpc/client:go_default_library",
        "//app/service/main/vip/api:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//app/service/main/vip/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
