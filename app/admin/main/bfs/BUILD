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
        "//app/admin/main/bfs/cmd:all-srcs",
        "//app/admin/main/bfs/conf:all-srcs",
        "//app/admin/main/bfs/dao:all-srcs",
        "//app/admin/main/bfs/http:all-srcs",
        "//app/admin/main/bfs/model:all-srcs",
        "//app/admin/main/bfs/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
