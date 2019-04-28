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
        "//app/service/ep/footman/cmd:all-srcs",
        "//app/service/ep/footman/conf:all-srcs",
        "//app/service/ep/footman/dao:all-srcs",
        "//app/service/ep/footman/model:all-srcs",
        "//app/service/ep/footman/server/http:all-srcs",
        "//app/service/ep/footman/service:all-srcs",
    ],
    tags = ["automanaged"],
)
