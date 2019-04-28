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
        "//app/interface/main/dm2/cmd:all-srcs",
        "//app/interface/main/dm2/conf:all-srcs",
        "//app/interface/main/dm2/dao:all-srcs",
        "//app/interface/main/dm2/http:all-srcs",
        "//app/interface/main/dm2/lib/xregex:all-srcs",
        "//app/interface/main/dm2/model:all-srcs",
        "//app/interface/main/dm2/rpc/client:all-srcs",
        "//app/interface/main/dm2/rpc/server:all-srcs",
        "//app/interface/main/dm2/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
