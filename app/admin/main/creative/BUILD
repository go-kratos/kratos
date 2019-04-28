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
        "//app/admin/main/creative/cmd:all-srcs",
        "//app/admin/main/creative/conf:all-srcs",
        "//app/admin/main/creative/dao:all-srcs",
        "//app/admin/main/creative/http:all-srcs",
        "//app/admin/main/creative/model/academy:all-srcs",
        "//app/admin/main/creative/model/app:all-srcs",
        "//app/admin/main/creative/model/bfs:all-srcs",
        "//app/admin/main/creative/model/logcli:all-srcs",
        "//app/admin/main/creative/model/material:all-srcs",
        "//app/admin/main/creative/model/music:all-srcs",
        "//app/admin/main/creative/model/operation:all-srcs",
        "//app/admin/main/creative/model/task:all-srcs",
        "//app/admin/main/creative/model/whitelist:all-srcs",
        "//app/admin/main/creative/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
