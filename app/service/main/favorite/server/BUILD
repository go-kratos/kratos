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
        "//app/service/main/favorite/server/gorpc:all-srcs",
        "//app/service/main/favorite/server/grpc:all-srcs",
        "//app/service/main/favorite/server/http:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
