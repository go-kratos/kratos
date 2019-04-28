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
        "//app/admin/ep/tapd/cmd:all-srcs",
        "//app/admin/ep/tapd/conf:all-srcs",
        "//app/admin/ep/tapd/dao:all-srcs",
        "//app/admin/ep/tapd/http:all-srcs",
        "//app/admin/ep/tapd/model:all-srcs",
        "//app/admin/ep/tapd/service:all-srcs",
    ],
    tags = ["automanaged"],
)
