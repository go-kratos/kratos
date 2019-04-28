package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dataflow.go",
        "display.go",
        "favorite.go",
        "history.go",
        "http.go",
        "local.go",
        "relation.go",
        "search.go",
        "space.go",
    ],
    importpath = "go-common/app/interface/main/app-interface/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/model:go_default_library",
        "//app/interface/main/app-interface/model/favorite:go_default_library",
        "//app/interface/main/app-interface/model/history:go_default_library",
        "//app/interface/main/app-interface/model/relation:go_default_library",
        "//app/interface/main/app-interface/model/search:go_default_library",
        "//app/interface/main/app-interface/model/space:go_default_library",
        "//app/interface/main/app-interface/service/account:go_default_library",
        "//app/interface/main/app-interface/service/dataflow:go_default_library",
        "//app/interface/main/app-interface/service/display:go_default_library",
        "//app/interface/main/app-interface/service/favorite:go_default_library",
        "//app/interface/main/app-interface/service/history:go_default_library",
        "//app/interface/main/app-interface/service/relation:go_default_library",
        "//app/interface/main/app-interface/service/search:go_default_library",
        "//app/interface/main/app-interface/service/space:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/proxy:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/xstr:go_default_library",
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
