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
        "//app/admin/ep/merlin/cmd:all-srcs",
        "//app/admin/ep/merlin/conf:all-srcs",
        "//app/admin/ep/merlin/dao:all-srcs",
        "//app/admin/ep/merlin/http:all-srcs",
        "//app/admin/ep/merlin/model:all-srcs",
        "//app/admin/ep/merlin/service:all-srcs",
    ],
    tags = ["automanaged"],
)
