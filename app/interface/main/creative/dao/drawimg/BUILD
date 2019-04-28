package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "draw.go",
        "freetype.go",
        "gaussian.go",
        "imgutil.go",
    ],
    importpath = "go-common/app/interface/main/creative/dao/drawimg",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/golang/freetype/raster:go_default_library",
        "//vendor/github.com/golang/freetype/truetype:go_default_library",
        "//vendor/golang.org/x/image/font:go_default_library",
        "//vendor/golang.org/x/image/math/fixed:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "draw_test.go",
        "freetype_test.go",
        "imgutil_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/golang/freetype/raster:go_default_library",
        "//vendor/github.com/golang/freetype/truetype:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/golang.org/x/image/font:go_default_library",
        "//vendor/golang.org/x/image/math/fixed:go_default_library",
    ],
)
