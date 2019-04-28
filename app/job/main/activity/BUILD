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
        "//app/job/main/activity/cmd:all-srcs",
        "//app/job/main/activity/conf:all-srcs",
        "//app/job/main/activity/dao/bnj:all-srcs",
        "//app/job/main/activity/dao/dm:all-srcs",
        "//app/job/main/activity/dao/kfc:all-srcs",
        "//app/job/main/activity/dao/like:all-srcs",
        "//app/job/main/activity/http:all-srcs",
        "//app/job/main/activity/model/bnj:all-srcs",
        "//app/job/main/activity/model/dm:all-srcs",
        "//app/job/main/activity/model/kfc:all-srcs",
        "//app/job/main/activity/model/like:all-srcs",
        "//app/job/main/activity/model/match:all-srcs",
        "//app/job/main/activity/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
