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
        "//app/job/main/coin/cmd:all-srcs",
        "//app/job/main/coin/conf:all-srcs",
        "//app/job/main/coin/dao:all-srcs",
        "//app/job/main/coin/http:all-srcs",
        "//app/job/main/coin/model:all-srcs",
        "//app/job/main/coin/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
