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
        "//app/job/main/app/cmd:all-srcs",
        "//app/job/main/app/conf:all-srcs",
        "//app/job/main/app/dao/feed:all-srcs",
        "//app/job/main/app/dao/monitor:all-srcs",
        "//app/job/main/app/dao/show:all-srcs",
        "//app/job/main/app/dao/space:all-srcs",
        "//app/job/main/app/dao/view:all-srcs",
        "//app/job/main/app/http:all-srcs",
        "//app/job/main/app/model:all-srcs",
        "//app/job/main/app/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
