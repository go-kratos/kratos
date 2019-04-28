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
        "//app/job/main/credit/cmd:all-srcs",
        "//app/job/main/credit/conf:all-srcs",
        "//app/job/main/credit/dao:all-srcs",
        "//app/job/main/credit/http:all-srcs",
        "//app/job/main/credit/model:all-srcs",
        "//app/job/main/credit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
