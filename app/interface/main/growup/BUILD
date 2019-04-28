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
        "//app/interface/main/growup/cmd:all-srcs",
        "//app/interface/main/growup/conf:all-srcs",
        "//app/interface/main/growup/dao:all-srcs",
        "//app/interface/main/growup/http:all-srcs",
        "//app/interface/main/growup/model:all-srcs",
        "//app/interface/main/growup/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
