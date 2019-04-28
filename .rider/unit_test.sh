#!/bin/bash

declare -a dirs=(dao)
declare -a packages
declare -a projects
declare mergeUser
declare commitUser
declare mergeID

#init env
function Init(){
    if [ ! -d "${CI_PROJECT_DIR}/../src" ];then
        mkdir ${CI_PROJECT_DIR}/../src
    fi
    ln -fs ${CI_PROJECT_DIR} ${CI_PROJECT_DIR}/../src
    export GOPATH=${CI_PROJECT_DIR}/..
}

function GetPackages(){
    reg="library/(.*)/(.*).go"
    for dir in ${dirs[@]}
    do
        reg+="|app/(service|interface|admin|job)/main/(.*)/${dir}/(.*).go"
    done
    files=`.rider/changefiles.sh ${suffix} | grep  -E "${reg}"`
    if [[ "${files}" = "" ]]; then
        echo "shell.GetPackages: no change files"
        exit 0
    fi
    for file in ${files}
    do
        # if [[ "${file}" =~ "library"* || "${file}" =~ "/mock" ]]; then
        if [[ "${file}" =~ "/mock" ]]; then
            continue
        fi
        package="go-common/$(dirname ${file})"
        if [[ ${packages} =~ ${package} ]]; then
            continue
        fi
        packages+=${package}" "
        packageSplit=(${package//\// })
        project=$(printf "%s/" ${packageSplit[@]:0:6})
        if [[ ${projects} =~ ${project} || ${project} =~ "/library" ]]; then
            continue
        fi
        projects+=${project%*/}" "
    done
    if [[ ${packages} = "" || ${projects} = "" ]]; then
        echo "shell.GetPackages no change packages"
        exit 0
    fi
}
# GetUserInfo get userinfo by gitlab result.
function GetUserInfo(){
    gitMergeRequestUrl="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/repository/commits/${CI_COMMIT_SHA}/merge_requests?private_token=${CI_PRIVATE_TOKEN}"
    gitCommitUrl="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/repository/commits/${CI_COMMIT_SHA}/statuses?private_token=${CI_PRIVATE_TOKEN}"
    mergeJson=$(curl -s ${gitMergeRequestUrl})
    commitJson=$(curl -s ${gitCommitUrl})
    if [[ "${mergeJson}" = "[]" ]] || [[ "${commitJson}" = "[]" ]]; then
        echo "Test not run, maybe you should try create a merge request first!"
        exit 0       
    fi
    mergeID=$(echo ${mergeJson} | jq -r '.[0].iid')
    mergeUser=$(echo ${mergeJson} | jq -r '.[0].author.username')
    commitUser=$(echo ${commitJson} | jq -r '.[0].author.username')
}

# Magic ignore method Check()
function Magic(){
    url="http://git.bilibili.co/api/v4/projects/${CI_PROJECT_ID}/merge_requests/${mergeID}/notes?private_token=${CI_PRIVATE_TOKEN}"
    json=$(curl -s ${url})
    admin="haoguanwei,chenjianrong,hedan,fengshanshan,zhaobingqing"
    len=$(echo ${json} | jq '.|length')
    for i in $(seq 0 $len)
    do
        comment=$(echo ${json} | jq -r ".[$i].body")
        user=$(echo ${json} | jq -r ".[$i].author.username")
        if [[ ${comment} = "+skiput" && ${admin} =~ ${user} ]]; then
             exit 0
        fi
    done
}

# Check determine whether the standard is up to standard
#$1: commit_id
function Check(){
    curl -s "${CI_UATSVEN_URL}/x/admin/apm/ut/git/report?project_id=${CI_PROJECT_ID}&merge_id=${mergeID}&commit_id=${CI_COMMIT_SHA}"
    checkURL="${CI_UATSVEN_URL}/x/admin/apm/ut/check?commit_id=${CI_COMMIT_SHA}"
    json=$(curl -s ${checkURL})
    code=$(echo ${json} | jq -r '.code')
    if [[ ${code} -ne 0 ]]; then
        echo -e "curl ${checkURL} response(${json})"
        exit 1
    fi
    package=$(echo ${json} | jq -r '.data.package')
    coverage=$(echo ${json} | jq -r '.data.coverage')
    passRate=$(echo ${json} | jq -r '.data.pass_rate')
    standard=$(echo ${json} | jq -r '.data.standard')
    increase=$(echo ${json} | jq -r '.data.increase')
    tyrant=$(echo ${json} | jq -r '.data.tyrant')
    lastCID=$(echo ${json} | jq -r '.data.last_cid')
    if ${tyrant}; then
        echo -e "\t续命失败!\n\t大佬，本次执行结果未达标哦(灬ꈍ ꈍ灬)，请再次优化ut重新提交🆙"
        echo -e "\t---------------------------------------------------------------------"
        printf "\t%-14s %-14s %-14s %-14s\n" "本次覆盖率(%)" "本次通过率(%)" "本次增长量(%)" 执行pkg
        printf "\t%-13.2f %-13.2f %-13.2f %-12s\n" ${coverage} ${passRate} ${increase} ${package}
        echo -e "\t(达标标准：覆盖率>=${standard} && 通过率=100% && 同比当前package历史最高覆盖率的增长率>=0)"
        echo -e "\t---------------------------------------------------------------------"
        echo -e "本次执行详细结果查询地址请访问：http://sven.bilibili.co/#/ut?merge_id=${mergeID}&&pn=1&ps=20"
        exit 1
    else
        echo -e "\t恭喜你，续命成功，可以请求合并MR了!"
    fi
}

function ReadDir(){
    # get go-common/app all dir path
    gopath=${GOPATH%..}
    PathDirs=`find ${gopath}app -maxdepth 3 -type d`
    value=""
    for dir in ${PathDirs}
    do
        if [[ -d "$dir" ]];then
            for file in `find ${dir} -maxdepth 1 -type f |grep "OWNERS"`
            do
                owner=""
                substr=${dir#*"go-common"}
                while read line
                do
                    if [[ "${line}" = "#"* ]] || [[ "${line}" = "" ]] || [[ "${line}" = "approvers:" ]];then
                        continue
                    elif [[ "${line}" = "labels:"* ]];then
                        break
                    else
                        owner+="${line:1},"
                    fi
                done < ${file}
                value+="{\"path\":\"go-common${substr}\",\"owner\":\"${owner%,}\"},"
            done
        fi
    done
    # delete "," at the end of value
    value=${value%,}
    echo "[${value}]" > path.out
}

# UTLint check the *_test.go files in the pkg
# $1: pkg
function UTLint()
{   
    cd $GOPATH
    declare -i numCase=0
    declare -i numAssertion=0
    files=$(ls $1 | grep -E "(.*)_test\.go")
    if [[ ${#files} -eq 0 ]];then
        echo "RunPKGUT.UTLint no *_test.go files in pkg($1)"
        exit 1
    fi
    for file in ${files}
    do
        numCase+=`grep -c -E "^func Test(.+)\(t \*testing\.T\) \{$" $1/${file}`
        numAssertion+=`grep -c -E "^(.*)So\((.+)\)$" $1/${file}`
    done
    if [[ ${numCase} -eq 0 || ${numAssertion} -eq 0 ]];then
        echo -e "RunPKGUT.UTLint no test case or assertion in pkg($1)"
        exit 1
    fi
}

# BazelTest execute bazel coverage and go tool
# $1: pkg
function BazelTest(){
    cd $GOPATH/go-common
    pkg=${1//go-common//}":go_default_test"
    path=${1//go-common\//}

    bazel coverage --config=ci --instrumentation_filter="//${path}[:],-//${path}/mock[/:]" --test_env=DEPLOY_ENV=uat --test_timeout=60 --test_env=APP_ID=bazel.test --test_output=all --cache_test_results=auto --test_arg=-convey-json ${pkg} > result.out
    if [[ ! -s result.out ]]; then 
        echo "==================================WARNING!======================================"
        echo "No test case found,请完善如下路径测试用例： ${pkg} "
        exit 1
    else
        echo $?
        cp $GOPATH/go-common/bazel-out/k8-fastbuild/testlogs/${path}/go_default_test/coverage.dat ./
        go tool cover -html=coverage.dat -o cover.html
    fi
}

# BazelTest execute bazel coverage for All files
# $1: pkg(go-common/app/admin/main/xxx/dao)
function BazelTestAll(){
    cd $GOPATH/go-common
    pkg=${1//go-common//}"/..."
    path=${1//go-common\//}
    echo "RunProjUT.BazelTestAll(${1}) pkg(${pkg}) path(${path}) pwd($(pwd))"
    bazel coverage --config=ci --instrumentation_filter="//${path}[/:],-//${path}/mock[/:]" --test_env=DEPLOY_ENV=uat --test_timeout=60 --test_env=APP_ID=bazel.test --test_output=all --cache_test_results=auto --test_arg=-convey-json ${pkg} > result.out
    find bazel-out/k8-fastbuild/testlogs/${path} -name "coverage.dat" | xargs cat | sort -nr | rev | uniq -s 1 | rev > coverage.dat
    coverage=$(cat coverage.dat | awk '{sum += $2;covSum += $2 * $3;} END {print covSum/sum*100}')
    sed -if "s/coverage: .*%/coverage: ${coverage}%/g" result.out
    go tool cover -html=coverage.dat -o cover.html
}

# upload data to apm
# $1: file result.out path
function Upload () {
    if [[ ! -f "result.out" ]] || [[ ! -f "cover.html" ]] || [[ ! -f "coverage.dat" ]]; then
        echo "==================================WARNING!======================================"
        echo "No test found!~ 请完善如下路径测试用例： ${1} "
        exit 1
    fi
    url="${CI_UATSVEN_URL}/x/admin/apm/ut/upload?merge_id=${mergeID}&username=${mergeUser}&author=${commitUser}&commit_id=${CI_COMMIT_SHA}&pkg=${1}"
    json=$(curl -s ${url} -H "Content-type: multipart/form-data" -F "html_file=@cover.html" -F "report_file=@result.out" -F "data_file=@coverage.dat")
    if [[ "${json}" = "" ]]; then
        echo "RunPKGUT.Upload curl ${url} fail"
        exit 1
    fi
    msg=$(echo ${json} | jq -r '.message')
    data=$(echo ${json} | jq -r '.data')
    code=$(echo ${json} | jq -r '.code')
    if [[ ${code} -ne 0 ]]; then
        echo "=============================================================================="
        echo -e "RunPKGUT.Upload Response. message(${msg})"
        echo -e "RunPKGUT.Upload Response. data(${data})\n\n"
        echo -e "RunPKGUT.Upload Upload Fail! status(${code})"
        exit ${code}
    fi
}

function RunPKGUT(){
    for package in ${packages}
    do
        echo "RunPKGUT.UTLint Start. pkg(${package})"
        UTLint ${package}
        echo "RunPKGUT.BazelTest Start. pkg(${package})"
        BazelTest ${package}
        echo "RunPKGUT.Upload Start. pkg(${package})"
        Upload ${package}
    done
    return 0
}

function RunProjUT(){
    for project in ${projects}
    do
        echo "RunProjUT.BazelTestAll Start. project(${project})"
        BazelTestAll ${project}
        echo "RunProjUT.Upload BazelTest Start. project(${project})"
        Upload ${project%/*}
    done
}

function UploadApp(){
    ReadDir
    url="${CI_UATSVEN_URL}/x/admin/apm/ut/upload/app" 
    curl -s ${url} -H "Content-type: multipart/form-data" -F "path_file=@path.out" > /dev/null
    echo "UploadApp() UpPath has finshed."
}

# run
Init
GetPackages
GetUserInfo
Magic
RunPKGUT
RunProjUT
UploadApp
Check
