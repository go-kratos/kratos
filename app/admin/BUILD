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
        "//app/admin/bbq:all-srcs",
        "//app/admin/ep:all-srcs",
        "//app/admin/live/live-admin:all-srcs",
        "//app/admin/main:all-srcs",
        "//app/admin/openplatform:all-srcs",
    ],
    tags = ["automanaged"],
)
