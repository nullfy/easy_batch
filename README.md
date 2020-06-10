# easy_batch
鉴于后端微服务架构，动辄上百个Repo，很多时候做个需求就变成了代码五分钟，分支两小时，遂有了此轮子。

# 安装
```shell
go get -v github.com/mumusa/easy_batch.git
```

# 使用
最后的`[repo_path]`不填默认已当前目录，easy_batch会递归找所有本地仓库，并执行git相关操作
 ```shell
 easy_batch git pull [repo_path]

 easy_batch git checkout master [repo_path]
 
 ```
 
 # ToDo
 * 支持根据文件批量下载仓库
 * 支持批量执行shell
