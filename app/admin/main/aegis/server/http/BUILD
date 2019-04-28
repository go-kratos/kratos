package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "csv.go",
        "debug.go",
        "direction.go",
        "engine.go",
        "flow.go",
        "http.go",
        "monitor.go",
        "net.go",
        "report.go",
        "svg.go",
        "task.go",
        "token.go",
        "transition.go",
    ],
    importpath = "go-common/app/admin/main/aegis/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/aegis/conf:go_default_library",
        "//app/admin/main/aegis/model:go_default_library",
        "//app/admin/main/aegis/model/business:go_default_library",
        "//app/admin/main/aegis/model/common:go_default_library",
        "//app/admin/main/aegis/model/middleware:go_default_library",
        "//app/admin/main/aegis/model/monitor:go_default_library",
        "//app/admin/main/aegis/model/net:go_default_library",
        "//app/admin/main/aegis/model/svg:go_default_library",
        "//app/admin/main/aegis/model/task:go_default_library",
        "//app/admin/main/aegis/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
