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
        "//app/job/main/passport/cmd:all-srcs",
        "//app/job/main/passport/conf:all-srcs",
        "//app/job/main/passport/dao:all-srcs",
        "//app/job/main/passport/http:all-srcs",
        "//app/job/main/passport/model:all-srcs",
        "//app/job/main/passport/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
