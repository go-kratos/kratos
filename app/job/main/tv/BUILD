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
        "//app/job/main/tv/cmd:all-srcs",
        "//app/job/main/tv/conf:all-srcs",
        "//app/job/main/tv/dao:all-srcs",
        "//app/job/main/tv/http:all-srcs",
        "//app/job/main/tv/model:all-srcs",
        "//app/job/main/tv/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
