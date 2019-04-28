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
        "//app/service/main/favorite/api:all-srcs",
        "//app/service/main/favorite/cmd:all-srcs",
        "//app/service/main/favorite/conf:all-srcs",
        "//app/service/main/favorite/dao:all-srcs",
        "//app/service/main/favorite/model:all-srcs",
        "//app/service/main/favorite/server:all-srcs",
        "//app/service/main/favorite/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
