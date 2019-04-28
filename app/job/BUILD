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
        "//app/job/bbq:all-srcs",
        "//app/job/live:all-srcs",
        "//app/job/live-userexp:all-srcs",
        "//app/job/main:all-srcs",
        "//app/job/openplatform:all-srcs",
    ],
    tags = ["automanaged"],
)
