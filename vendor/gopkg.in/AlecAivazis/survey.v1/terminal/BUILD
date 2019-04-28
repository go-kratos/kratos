package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "buffered_reader.go",
        "cursor.go",
        "cursor_windows.go",
        "display.go",
        "display_posix.go",
        "display_windows.go",
        "error.go",
        "output.go",
        "output_windows.go",
        "runereader.go",
        "runereader_bsd.go",
        "runereader_linux.go",
        "runereader_posix.go",
        "runereader_windows.go",
        "sequences.go",
        "stdio.go",
        "syscall_windows.go",
        "terminal.go",
    ],
    importpath = "gopkg.in/AlecAivazis/survey.v1/terminal",
    tags = ["manual"],
    visibility = ["//visibility:public"],
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
