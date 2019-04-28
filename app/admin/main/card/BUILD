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
        "//app/admin/main/card/cmd:all-srcs",
        "//app/admin/main/card/conf:all-srcs",
        "//app/admin/main/card/dao:all-srcs",
        "//app/admin/main/card/model:all-srcs",
        "//app/admin/main/card/server/http:all-srcs",
        "//app/admin/main/card/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
