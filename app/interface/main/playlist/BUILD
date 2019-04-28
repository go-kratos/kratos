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
        "//app/interface/main/playlist/cmd:all-srcs",
        "//app/interface/main/playlist/conf:all-srcs",
        "//app/interface/main/playlist/dao:all-srcs",
        "//app/interface/main/playlist/http:all-srcs",
        "//app/interface/main/playlist/model:all-srcs",
        "//app/interface/main/playlist/rpc/client:all-srcs",
        "//app/interface/main/playlist/rpc/server:all-srcs",
        "//app/interface/main/playlist/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
