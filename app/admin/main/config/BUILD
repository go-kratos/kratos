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
        "//app/admin/main/config/cmd:all-srcs",
        "//app/admin/main/config/conf:all-srcs",
        "//app/admin/main/config/dao:all-srcs",
        "//app/admin/main/config/http:all-srcs",
        "//app/admin/main/config/model:all-srcs",
        "//app/admin/main/config/pkg/lint:all-srcs",
        "//app/admin/main/config/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
