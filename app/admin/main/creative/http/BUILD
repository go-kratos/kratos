package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "academy_test.go",
        "whiteList_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/creative/conf:go_default_library",
        "//app/admin/main/creative/service:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "academy_archive.go",
        "academy_occupation.go",
        "academy_tag.go",
        "app.go",
        "csv.go",
        "export.go",
        "http.go",
        "index.go",
        "local.go",
        "material.go",
        "material_category.go",
        "music.go",
        "music_category.go",
        "music_material.go",
        "music_with_category.go",
        "music_with_material.go",
        "operCollectArc.go",
        "task.go",
        "upload.go",
        "whitelist.go",
    ],
    importpath = "go-common/app/admin/main/creative/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/creative/conf:go_default_library",
        "//app/admin/main/creative/model/academy:go_default_library",
        "//app/admin/main/creative/model/app:go_default_library",
        "//app/admin/main/creative/model/logcli:go_default_library",
        "//app/admin/main/creative/model/material:go_default_library",
        "//app/admin/main/creative/model/music:go_default_library",
        "//app/admin/main/creative/model/operation:go_default_library",
        "//app/admin/main/creative/model/task:go_default_library",
        "//app/admin/main/creative/model/whitelist:go_default_library",
        "//app/admin/main/creative/service:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
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
