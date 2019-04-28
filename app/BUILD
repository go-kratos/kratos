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
        "//app/admin:all-srcs",
        "//app/common:all-srcs",
        "//app/infra:all-srcs",
        "//app/interface:all-srcs",
        "//app/job:all-srcs",
        "//app/service:all-srcs",
        "//app/tool:all-srcs",
    ],
    tags = ["automanaged"],
)
