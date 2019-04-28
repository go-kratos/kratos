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
        "//app/admin/main/feed/cmd:all-srcs",
        "//app/admin/main/feed/conf:all-srcs",
        "//app/admin/main/feed/dao/account:all-srcs",
        "//app/admin/main/feed/dao/archive:all-srcs",
        "//app/admin/main/feed/dao/bfs:all-srcs",
        "//app/admin/main/feed/dao/egg:all-srcs",
        "//app/admin/main/feed/dao/pgc:all-srcs",
        "//app/admin/main/feed/dao/search:all-srcs",
        "//app/admin/main/feed/dao/show:all-srcs",
        "//app/admin/main/feed/http:all-srcs",
        "//app/admin/main/feed/model:all-srcs",
        "//app/admin/main/feed/service:all-srcs",
        "//app/admin/main/feed/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
