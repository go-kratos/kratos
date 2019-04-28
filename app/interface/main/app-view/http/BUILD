package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bnj.go",
        "http.go",
        "local.go",
        "report.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/app-view/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-view/conf:go_default_library",
        "//app/interface/main/app-view/model:go_default_library",
        "//app/interface/main/app-view/model/view:go_default_library",
        "//app/interface/main/app-view/service/report:go_default_library",
        "//app/interface/main/app-view/service/view:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/anticheat:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/proxy:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
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
