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
        "//app/job/main/creative/cmd:all-srcs",
        "//app/job/main/creative/conf:all-srcs",
        "//app/job/main/creative/dao/academy:all-srcs",
        "//app/job/main/creative/dao/archive:all-srcs",
        "//app/job/main/creative/dao/monitor:all-srcs",
        "//app/job/main/creative/dao/newcomer:all-srcs",
        "//app/job/main/creative/dao/weeklyhonor:all-srcs",
        "//app/job/main/creative/http:all-srcs",
        "//app/job/main/creative/model:all-srcs",
        "//app/job/main/creative/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
