package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cookie.go",
        "cookie_go111.go",
        "doc.go",
        "lex.go",
        "options.go",
        "options_go111.go",
        "sessions.go",
        "store.go",
    ],
    importpath = "github.com/gorilla/sessions",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/gorilla/context:go_default_library",
        "//vendor/github.com/gorilla/securecookie:go_default_library",
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
