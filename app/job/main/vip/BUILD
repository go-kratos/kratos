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
        "//app/job/main/vip/cmd:all-srcs",
        "//app/job/main/vip/conf:all-srcs",
        "//app/job/main/vip/dao:all-srcs",
        "//app/job/main/vip/http:all-srcs",
        "//app/job/main/vip/model:all-srcs",
        "//app/job/main/vip/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
