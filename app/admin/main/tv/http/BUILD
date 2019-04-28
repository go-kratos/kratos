package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "arc_audit.go",
        "archive.go",
        "audit_result.go",
        "channel.go",
        "content_repo.go",
        "ep_content.go",
        "http.go",
        "intervs.go",
        "label.go",
        "mango.go",
        "modules.go",
        "order.go",
        "panel.go",
        "region.go",
        "sear_inter.go",
        "season_repo.go",
        "upbfs.go",
        "upper.go",
        "user.go",
        "version.go",
        "version_update.go",
        "video.go",
        "watermark.go",
    ],
    importpath = "go-common/app/admin/main/tv/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/tv/conf:go_default_library",
        "//app/admin/main/tv/model:go_default_library",
        "//app/admin/main/tv/service:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
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
