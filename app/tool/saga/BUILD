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
        "//app/tool/saga/cmd:all-srcs",
        "//app/tool/saga/conf:all-srcs",
        "//app/tool/saga/dao:all-srcs",
        "//app/tool/saga/http:all-srcs",
        "//app/tool/saga/model:all-srcs",
        "//app/tool/saga/service:all-srcs",
    ],
    tags = ["automanaged"],
)
