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
        "//app/job/live-userexp/cmd:all-srcs",
        "//app/job/live-userexp/conf:all-srcs",
        "//app/job/live-userexp/dao:all-srcs",
        "//app/job/live-userexp/model:all-srcs",
        "//app/job/live-userexp/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
