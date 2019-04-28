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
        "//app/admin/main/mcn/cmd:all-srcs",
        "//app/admin/main/mcn/conf:all-srcs",
        "//app/admin/main/mcn/dao/bfs:all-srcs",
        "//app/admin/main/mcn/dao/msg:all-srcs",
        "//app/admin/main/mcn/dao/up:all-srcs",
        "//app/admin/main/mcn/dao/videoup:all-srcs",
        "//app/admin/main/mcn/model:all-srcs",
        "//app/admin/main/mcn/server/http:all-srcs",
        "//app/admin/main/mcn/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
