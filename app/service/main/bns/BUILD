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
        "//app/service/main/bns/agent:all-srcs",
        "//app/service/main/bns/cmd:all-srcs",
        "//app/service/main/bns/conf:all-srcs",
        "//app/service/main/bns/lib:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
