load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "comment.go",
        "data.go",
        "feed.go",
        "http.go",
        "location.go",
        "notice.go",
        "push.go",
        "report.go",
        "search.go",
        "share.go",
        "sv.go",
        "topic.go",
        "upload.go",
        "user.go",
    ],
    importpath = "go-common/app/interface/bbq/app-bbq/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/bbq/app-bbq/api/http/v1:go_default_library",
        "//app/interface/bbq/app-bbq/conf:go_default_library",
        "//app/interface/bbq/app-bbq/model:go_default_library",
        "//app/interface/bbq/app-bbq/service:go_default_library",
        "//app/interface/bbq/common/auth:go_default_library",
        "//app/interface/bbq/common/http:go_default_library",
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/trace:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
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
