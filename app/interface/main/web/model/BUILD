package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "abserver.go",
        "account.go",
        "archive.go",
        "article.go",
        "baidu.go",
        "bigdata.go",
        "bnj.go",
        "coin.go",
        "elec.go",
        "feedback.go",
        "help.go",
        "icon.go",
        "nav.go",
        "online.go",
        "reply.go",
        "search.go",
        "shop.go",
        "tag.go",
        "web.go",
        "wechat.go",
    ],
    importpath = "go-common/app/interface/main/web/model",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/dm2/model:go_default_library",
        "//app/interface/main/reply/model/reply:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/broadcast/api/grpc/v1:go_default_library",
        "//app/service/main/ugcpay/api/grpc/v1:go_default_library",
        "//library/time:go_default_library",
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
