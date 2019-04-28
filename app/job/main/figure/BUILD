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
        "//app/job/main/figure/cmd:all-srcs",
        "//app/job/main/figure/conf:all-srcs",
        "//app/job/main/figure/dao:all-srcs",
        "//app/job/main/figure/http:all-srcs",
        "//app/job/main/figure/model:all-srcs",
        "//app/job/main/figure/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
