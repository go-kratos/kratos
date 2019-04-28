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
        "//app/interface/main/reply/cmd:all-srcs",
        "//app/interface/main/reply/conf:all-srcs",
        "//app/interface/main/reply/dao/bigdata:all-srcs",
        "//app/interface/main/reply/dao/drawyoo:all-srcs",
        "//app/interface/main/reply/dao/fans:all-srcs",
        "//app/interface/main/reply/dao/reply:all-srcs",
        "//app/interface/main/reply/dao/search:all-srcs",
        "//app/interface/main/reply/dao/vip:all-srcs",
        "//app/interface/main/reply/dao/workflow:all-srcs",
        "//app/interface/main/reply/http:all-srcs",
        "//app/interface/main/reply/model/activity:all-srcs",
        "//app/interface/main/reply/model/adminlog:all-srcs",
        "//app/interface/main/reply/model/drawyoo:all-srcs",
        "//app/interface/main/reply/model/reply:all-srcs",
        "//app/interface/main/reply/model/topic:all-srcs",
        "//app/interface/main/reply/model/vip:all-srcs",
        "//app/interface/main/reply/model/xreply:all-srcs",
        "//app/interface/main/reply/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
