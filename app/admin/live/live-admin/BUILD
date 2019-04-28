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
        "//app/admin/live/live-admin/api/http/v1:all-srcs",
        "//app/admin/live/live-admin/api/http/v2:all-srcs",
        "//app/admin/live/live-admin/cmd:all-srcs",
        "//app/admin/live/live-admin/conf:all-srcs",
        "//app/admin/live/live-admin/dao:all-srcs",
        "//app/admin/live/live-admin/model:all-srcs",
        "//app/admin/live/live-admin/server/http:all-srcs",
        "//app/admin/live/live-admin/service:all-srcs",
    ],
    tags = ["automanaged"],
)
