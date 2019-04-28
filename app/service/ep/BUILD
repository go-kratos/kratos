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
        "//app/service/ep/footman:all-srcs",
        "//app/service/ep/saga-agent:all-srcs",
    ],
    tags = ["automanaged"],
)
