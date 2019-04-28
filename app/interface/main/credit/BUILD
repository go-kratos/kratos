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
        "//app/interface/main/credit/cmd:all-srcs",
        "//app/interface/main/credit/conf:all-srcs",
        "//app/interface/main/credit/dao:all-srcs",
        "//app/interface/main/credit/http:all-srcs",
        "//app/interface/main/credit/model:all-srcs",
        "//app/interface/main/credit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
