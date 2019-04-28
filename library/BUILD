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
        "//library/cache:all-srcs",
        "//library/conf:all-srcs",
        "//library/container:all-srcs",
        "//library/database:all-srcs",
        "//library/ecode:all-srcs",
        "//library/exp/feature:all-srcs",
        "//library/log:all-srcs",
        "//library/naming:all-srcs",
        "//library/net:all-srcs",
        "//library/os:all-srcs",
        "//library/queue:all-srcs",
        "//library/rate:all-srcs",
        "//library/stat:all-srcs",
        "//library/sync:all-srcs",
        "//library/syscall:all-srcs",
        "//library/text/translate/chinese:all-srcs",
        "//library/time:all-srcs",
        "//library/xstr:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
