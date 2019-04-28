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
        "//app/admin/main/up-rating/cmd:all-srcs",
        "//app/admin/main/up-rating/conf:all-srcs",
        "//app/admin/main/up-rating/dao:all-srcs",
        "//app/admin/main/up-rating/http:all-srcs",
        "//app/admin/main/up-rating/model:all-srcs",
        "//app/admin/main/up-rating/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
