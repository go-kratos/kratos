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
        "//app/job/main/videoup/dao/activity:all-srcs",
        "//app/job/main/videoup/dao/archive:all-srcs",
        "//app/job/main/videoup/dao/bvc:all-srcs",
        "//app/job/main/videoup/dao/manager:all-srcs",
        "//app/job/main/videoup/dao/message:all-srcs",
        "//app/job/main/videoup/dao/monitor:all-srcs",
        "//app/job/main/videoup/dao/redis:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
