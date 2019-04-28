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
        "//app/interface/main/space/cmd:all-srcs",
        "//app/interface/main/space/conf:all-srcs",
        "//app/interface/main/space/dao:all-srcs",
        "//app/interface/main/space/http:all-srcs",
        "//app/interface/main/space/model:all-srcs",
        "//app/interface/main/space/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
