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
        "//app/admin/ep/marthe:all-srcs",
        "//app/admin/ep/melloi:all-srcs",
        "//app/admin/ep/merlin:all-srcs",
        "//app/admin/ep/saga:all-srcs",
        "//app/admin/ep/tapd:all-srcs",
    ],
    tags = ["automanaged"],
)
