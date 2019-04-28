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
        "//app/job/main/archive/cmd:all-srcs",
        "//app/job/main/archive/conf:all-srcs",
        "//app/job/main/archive/dao/archive:all-srcs",
        "//app/job/main/archive/dao/email:all-srcs",
        "//app/job/main/archive/dao/monitor:all-srcs",
        "//app/job/main/archive/dao/reply:all-srcs",
        "//app/job/main/archive/dao/result:all-srcs",
        "//app/job/main/archive/http:all-srcs",
        "//app/job/main/archive/model/archive:all-srcs",
        "//app/job/main/archive/model/databus:all-srcs",
        "//app/job/main/archive/model/dm:all-srcs",
        "//app/job/main/archive/model/pgc:all-srcs",
        "//app/job/main/archive/model/result:all-srcs",
        "//app/job/main/archive/model/retry:all-srcs",
        "//app/job/main/archive/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
