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
        "//app/job/main/stat/cmd:all-srcs",
        "//app/job/main/stat/conf:all-srcs",
        "//app/job/main/stat/dao:all-srcs",
        "//app/job/main/stat/http:all-srcs",
        "//app/job/main/stat/model:all-srcs",
        "//app/job/main/stat/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
