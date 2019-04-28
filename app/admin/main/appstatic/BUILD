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
        "//app/admin/main/appstatic/cmd:all-srcs",
        "//app/admin/main/appstatic/conf:all-srcs",
        "//app/admin/main/appstatic/dao:all-srcs",
        "//app/admin/main/appstatic/http:all-srcs",
        "//app/admin/main/appstatic/model:all-srcs",
        "//app/admin/main/appstatic/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
