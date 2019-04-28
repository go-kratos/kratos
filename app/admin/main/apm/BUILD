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
        "//app/admin/main/apm/cmd:all-srcs",
        "//app/admin/main/apm/conf:all-srcs",
        "//app/admin/main/apm/dao:all-srcs",
        "//app/admin/main/apm/http:all-srcs",
        "//app/admin/main/apm/model:all-srcs",
        "//app/admin/main/apm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
