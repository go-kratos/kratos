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
        "//app/service/main/account/rpc/client:all-srcs",
        "//app/service/main/account/rpc/server:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
