package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "cm.go",
        "flow.go",
        "http.go",
        "local.go",
        "monitor.go",
        "music.go",
        "oversea.go",
        "pgc.go",
        "search.go",
        "staff.go",
        "stats.go",
        "task.go",
        "track.go",
        "video.go",
    ],
    importpath = "go-common/app/admin/main/videoup/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup/conf:go_default_library",
        "//app/admin/main/videoup/model/archive:go_default_library",
        "//app/admin/main/videoup/model/monitor:go_default_library",
        "//app/admin/main/videoup/model/music:go_default_library",
        "//app/admin/main/videoup/model/oversea:go_default_library",
        "//app/admin/main/videoup/model/search:go_default_library",
        "//app/admin/main/videoup/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
