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
        "//app/admin/main/app/cmd:all-srcs",
        "//app/admin/main/app/conf:all-srcs",
        "//app/admin/main/app/dao/aids:all-srcs",
        "//app/admin/main/app/dao/audit:all-srcs",
        "//app/admin/main/app/dao/bfs:all-srcs",
        "//app/admin/main/app/dao/bottom:all-srcs",
        "//app/admin/main/app/dao/language:all-srcs",
        "//app/admin/main/app/dao/notice:all-srcs",
        "//app/admin/main/app/dao/wall:all-srcs",
        "//app/admin/main/app/http:all-srcs",
        "//app/admin/main/app/model:all-srcs",
        "//app/admin/main/app/service/aids:all-srcs",
        "//app/admin/main/app/service/audit:all-srcs",
        "//app/admin/main/app/service/bfs:all-srcs",
        "//app/admin/main/app/service/bottom:all-srcs",
        "//app/admin/main/app/service/language:all-srcs",
        "//app/admin/main/app/service/notice:all-srcs",
        "//app/admin/main/app/service/ping:all-srcs",
        "//app/admin/main/app/service/wall:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
