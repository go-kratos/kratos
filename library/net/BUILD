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
        "//library/net/http:all-srcs",
        "//library/net/ip:all-srcs",
        "//library/net/metadata:all-srcs",
        "//library/net/netutil:all-srcs",
        "//library/net/rpc:all-srcs",
        "//library/net/trace:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
