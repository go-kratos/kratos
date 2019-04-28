package(default_visibility = ["//visibility:public"])

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/tool/ci/cmd/mail:all-srcs",
        "//app/tool/ci/lib/changes:all-srcs",
        "//app/tool/ci/lib/commit:all-srcs",
        "//app/tool/ci/lib/lint:all-srcs",
        "//app/tool/ci/lib/mail:all-srcs",
    ],
    tags = ["automanaged"],
)
