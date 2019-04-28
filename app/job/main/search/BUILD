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
        "//app/job/main/search/cmd:all-srcs",
        "//app/job/main/search/conf:all-srcs",
        "//app/job/main/search/dao:all-srcs",
        "//app/job/main/search/http:all-srcs",
        "//app/job/main/search/model:all-srcs",
        "//app/job/main/search/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
