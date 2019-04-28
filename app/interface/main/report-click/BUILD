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
        "//app/interface/main/report-click/cmd:all-srcs",
        "//app/interface/main/report-click/conf:all-srcs",
        "//app/interface/main/report-click/dao:all-srcs",
        "//app/interface/main/report-click/http:all-srcs",
        "//app/interface/main/report-click/model:all-srcs",
        "//app/interface/main/report-click/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
