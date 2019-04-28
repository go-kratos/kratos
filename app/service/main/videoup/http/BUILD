package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "flow.go",
        "http.go",
        "local.go",
        "netsafe.go",
        "pgc.go",
        "porder.go",
        "reco.go",
        "report.go",
        "staff.go",
        "ugc.go",
        "video.go",
    ],
    importpath = "go-common/app/service/main/videoup/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/videoup/conf:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
        "//app/service/main/videoup/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
