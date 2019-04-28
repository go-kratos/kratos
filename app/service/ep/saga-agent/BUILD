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
        "//app/service/ep/saga-agent/cmd:all-srcs",
        "//app/service/ep/saga-agent/conf:all-srcs",
        "//app/service/ep/saga-agent/path:all-srcs",
        "//app/service/ep/saga-agent/service/agent:all-srcs",
    ],
    tags = ["automanaged"],
)
