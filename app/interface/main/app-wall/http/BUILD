package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "ip.go",
        "local.go",
        "mobile.go",
        "offer.go",
        "operator.go",
        "telecom.go",
        "unicom.go",
        "wall.go",
    ],
    importpath = "go-common/app/interface/main/app-wall/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-wall/conf:go_default_library",
        "//app/interface/main/app-wall/dao/padding:go_default_library",
        "//app/interface/main/app-wall/model:go_default_library",
        "//app/interface/main/app-wall/model/mobile:go_default_library",
        "//app/interface/main/app-wall/model/telecom:go_default_library",
        "//app/interface/main/app-wall/model/unicom:go_default_library",
        "//app/interface/main/app-wall/service/mobile:go_default_library",
        "//app/interface/main/app-wall/service/offer:go_default_library",
        "//app/interface/main/app-wall/service/operator:go_default_library",
        "//app/interface/main/app-wall/service/ping:go_default_library",
        "//app/interface/main/app-wall/service/telecom:go_default_library",
        "//app/interface/main/app-wall/service/unicom:go_default_library",
        "//app/interface/main/app-wall/service/wall:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/proxy:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
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
