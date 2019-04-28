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
        "//app/admin/ep/melloi/cmd:all-srcs",
        "//app/admin/ep/melloi/conf:all-srcs",
        "//app/admin/ep/melloi/dao:all-srcs",
        "//app/admin/ep/melloi/http:all-srcs",
        "//app/admin/ep/melloi/model:all-srcs",
        "//app/admin/ep/melloi/service:all-srcs",
    ],
    tags = ["automanaged"],
)
