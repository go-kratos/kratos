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
        "//app/service/main/dynamic/cmd:all-srcs",
        "//app/service/main/dynamic/conf:all-srcs",
        "//app/service/main/dynamic/dao:all-srcs",
        "//app/service/main/dynamic/http:all-srcs",
        "//app/service/main/dynamic/model:all-srcs",
        "//app/service/main/dynamic/rpc:all-srcs",
        "//app/service/main/dynamic/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
