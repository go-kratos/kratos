package(default_visibility = ["//visibility:public"])

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
        "//app/admin/ep/marthe/cmd:all-srcs",
        "//app/admin/ep/marthe/conf:all-srcs",
        "//app/admin/ep/marthe/dao:all-srcs",
        "//app/admin/ep/marthe/model:all-srcs",
        "//app/admin/ep/marthe/server/http:all-srcs",
        "//app/admin/ep/marthe/service:all-srcs",
    ],
    tags = ["automanaged"],
)
