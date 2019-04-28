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
        "//app/job/main/figure-timer/cmd:all-srcs",
        "//app/job/main/figure-timer/conf:all-srcs",
        "//app/job/main/figure-timer/dao:all-srcs",
        "//app/job/main/figure-timer/http:all-srcs",
        "//app/job/main/figure-timer/model:all-srcs",
        "//app/job/main/figure-timer/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
