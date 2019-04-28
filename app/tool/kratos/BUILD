package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = [
        "build.go",
        "doc.go",
        "init.go",
        "main.go",
        "project.go",
        "template.go",
        "update.go",
        "upgrade.go",
        "verification.go",
        "version.go",
    ],
    importpath = "go-common/app/tool/kratos",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/AlecAivazis/survey:go_default_library",
        "//vendor/github.com/kr/pty:go_default_library",
        "//vendor/github.com/urfave/cli:go_default_library",
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

go_binary(
    name = "kratos",
    embed = [":go_default_library"],
    tags = ["automanaged"],
    x_defs = {
        "main.Version": "3",
        "main.GitCommit": "{STABLE_BUILD_GIT_COMMIT}",
        "main.BuildTime": "{STABLE_BUILD_TIME}",
        "main.Channel": "release",
    },
)
