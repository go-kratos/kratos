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
        "//app/job/main/reply/cmd:all-srcs",
        "//app/job/main/reply/conf:all-srcs",
        "//app/job/main/reply/dao/message:all-srcs",
        "//app/job/main/reply/dao/notice:all-srcs",
        "//app/job/main/reply/dao/reply:all-srcs",
        "//app/job/main/reply/dao/search:all-srcs",
        "//app/job/main/reply/dao/spam:all-srcs",
        "//app/job/main/reply/dao/stat:all-srcs",
        "//app/job/main/reply/http:all-srcs",
        "//app/job/main/reply/model/activity:all-srcs",
        "//app/job/main/reply/model/drawyoo:all-srcs",
        "//app/job/main/reply/model/reply:all-srcs",
        "//app/job/main/reply/model/topic:all-srcs",
        "//app/job/main/reply/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
