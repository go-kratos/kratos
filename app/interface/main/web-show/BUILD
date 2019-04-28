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
        "//app/interface/main/web-show/cmd:all-srcs",
        "//app/interface/main/web-show/conf:all-srcs",
        "//app/interface/main/web-show/dao/ad:all-srcs",
        "//app/interface/main/web-show/dao/bangumi:all-srcs",
        "//app/interface/main/web-show/dao/data:all-srcs",
        "//app/interface/main/web-show/dao/job:all-srcs",
        "//app/interface/main/web-show/dao/operation:all-srcs",
        "//app/interface/main/web-show/dao/resource:all-srcs",
        "//app/interface/main/web-show/http:all-srcs",
        "//app/interface/main/web-show/model/job:all-srcs",
        "//app/interface/main/web-show/model/operation:all-srcs",
        "//app/interface/main/web-show/model/resource:all-srcs",
        "//app/interface/main/web-show/service/job:all-srcs",
        "//app/interface/main/web-show/service/operation:all-srcs",
        "//app/interface/main/web-show/service/resource:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
