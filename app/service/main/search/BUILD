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
        "//app/service/main/search/cmd:all-srcs",
        "//app/service/main/search/conf:all-srcs",
        "//app/service/main/search/dao:all-srcs",
        "//app/service/main/search/http:all-srcs",
        "//app/service/main/search/model:all-srcs",
        "//app/service/main/search/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
