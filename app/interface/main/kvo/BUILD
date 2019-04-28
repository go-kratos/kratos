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
        "//app/interface/main/kvo/cmd:all-srcs",
        "//app/interface/main/kvo/conf:all-srcs",
        "//app/interface/main/kvo/dao:all-srcs",
        "//app/interface/main/kvo/http:all-srcs",
        "//app/interface/main/kvo/model:all-srcs",
        "//app/interface/main/kvo/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
