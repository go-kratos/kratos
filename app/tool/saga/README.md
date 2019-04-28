# saga

##### 项目简介
> 1.提供大仓库pkg依赖关系DAG  
> 2.提供gitlab MR自动构建、测试、覆盖率、代码静态检查  

##### 编译环境
> 请只用golang v1.8.x以上版本编译执行。  

##### 依赖包
> 1.公共仓库go-common  
> 2.github.com/xanzy/go-gitlab  

##### 依赖服务
> 1.gitlab

##### 特别说明
> 1.运行环境ssh key需要绑定到gitlab账户下
> 2.运行环境PATH有运行的go,golint的路径
> 3.eslint 
    1. 安装nodejs
    2. 设置path
    3. npm run lint
> 4.phplint
    1. 安装 PHP& PEAR
    2. 安装CodeSniffer
        pear install PHP_CodeSniffer
    3. 配置CodeSniffer
        phpcs --config-set default_standard PSR2
        phpcs --config-set show_warnings 0
        phpcs --config-set severity 1
    4. 校验代码
        phpcs * (*为待校验的文件名，可以只检验本次MR涉及改动的文件)
    