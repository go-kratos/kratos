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
        "//app/job/live/wallet/cmd:all-srcs",
        "//app/job/live/wallet/conf:all-srcs",
        "//app/job/live/wallet/dao:all-srcs",
        "//app/job/live/wallet/model:all-srcs",
        "//app/job/live/wallet/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
