package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "captcha_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/captcha/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "bilinear.go",
        "business.go",
        "captcha.go",
        "draw.go",
        "image.go",
        "rotate.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/captcha/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/captcha/conf:go_default_library",
        "//app/interface/main/captcha/dao:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/golang/freetype:go_default_library",
        "//vendor/github.com/golang/freetype/truetype:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
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
