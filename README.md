# 万能文章采集器(collector)
这是一个由golang编写的采集器，可以自动识别文章列表和文章内容。使用它来采集文章并不需要编写正则表达式，你只需要提供文章列表页的连接即可。

## 为什么会有这个万能文章采集器
* 市面上有几种采集工具，大多都需要针对不同的网站定制不同的采集规则，才能最终采集到想要的结果。本采集器内置了常用的采集规则，只要添加文章列表连接，就能将内容采集回来。
* 本采集器采用多线程并行采集，可在同一时间采集更多的内容。
* 本采集器只专注于采集文章这一件事，不用来定制采集其他内容，只专心做一件事。

## 万能文章采集器能采集哪些内容
本采集器可以采集到的的内容有：文章标题、文章关键词、文章描述、文章详情内容、文章作者、文章发布时间、文章浏览量。

## 什么时候需要使用到万能文章采集器
当我们需要给网站采集文章的时候，本采集器就可以派上用场了，本采集器不需要有人值守，24小时不间断运行，每隔10分钟就会自动遍历一遍采集列表，抓取包含有文章的连接，随时将文字抓取回来，还可以设置自动发布，自动发布到指定文章表中。

## 万能文章采集器可用在哪里运行
本采集器可用运行在 Windows系统、Mac 系统、Linux系统（Centos、Ubuntu等），可用下载编译好的程序直接执行，也可以下载源码自己编译。

## 万能文章采集器可用伪原创吗
本采集器暂时还不支持伪原创功能，后期会增加适当的伪原创选项。

## 如何安装使用
* 下载可执行文件  
  请从Releases 中根据你的操作系统下载最新版的可执行文件，解压后，重命名config.dist.json为config.json，打开config.json，修改mysql部分的配置，填写为你的mysql地址、用户名、密码、数据库信息，导入mysql.sql到填写的数据库中，然后双击运行可执行文件即可开始采集之旅。
* 自助编译  
  先clone代码到本地，本地安装go运行环境，在collector目录下打开cmd/Terminal命令行窗口，执行命。如果你没配置代理的话，还需要新设置go的代理
```shell script
go env -w GOPROXY=https://goproxy.cn,direct
```
  最后执行下面命令  
```shell script
go mod tidy
go mod vendor
go build
```
编译结束后，配置config。重命名config.dist.json为config.json，打开config.json，修改mysql部分的配置，填写为你的mysql地址、用户名、密码、数据库信息，导入mysql.sql到填写的数据库中，然后双击运行可执行文件即可开始采集之旅。

### 添加待采集文章列表说明
第一版尚未有可视化界面，因此需要你使用数据库工具打开fe_article_source 表，在里面填充采集列表，只需要将需要采集的列表填写到url字段即可，一行一个。

### config.json配置说明
```
{
  "mysql": { //数据库配置
    "Database": "collector",
    "User": "root",
    "Password": "root",
    "Charset": "utf8mb4",
    "Host": "127.0.0.1",
    "TablePrefix": "fe_",
    "Port": 3306,
    "MaxIdleConnections": 1000,
    "MaxOpenConnections": 100000
  },
  "server": { //采集器运行配置
    "SiteName"    : "万能采集器",
    "Host"        : "localhost",
    "Env"         : "development",
    "Port"        : 8088
  },
  "collector": { //采集规则
    "ErrorTimes": 5, //列表访问错误多少次后抛弃该列表连接
    "Channels": 5,  //同时使用多少个通道执行
    "TitleMinLength": 6,  //最小标题长度，小于该长度的会自动放弃
    "ContentMinLength": 200,  //最小详情长度，小于该长度的会自动放弃
    "TitleExclude": [  //标题不包含关键词，出现这些关键词的会自动放弃
      "法律声明",
      "关于我们",
      "站点地图"
    ],
    "TitleExcludePrefix": [  //标题不包含开头，以这些开头的会自动放弃
      "404",
      "403",
      "NotFound"
    ],
    "TitleExcludeSuffix": [  //标题不包含结尾，以这些开头的会自动放弃
      "网站",
      "网",
      "政府",
      "门户"
    ],
    "ContentExclude": [  //内容不包含关键词，出现这些关键词的会自动放弃
      "ICP备",
      "政府网站标识码",
      "以上版本浏览本站",
      "版权声明",
      "公网安备"
    ],
    "ContentExcludeLine": [  //内容不包含关键词的行，出现这些关键词的行会自动放弃
      "背景色：",
      "时间：",
      "作者：",
      "qrcode"
    ]
  },
  "content": {  //自动发布设置
    "AutoPublish": true,  //是否自动发布，true为自动
    "TableName": "fe_new_article",  //自动发布到的文章表名
    "IdField": "id",  //文章表的id字段名
    "TitleField": "title",  //文章表的标题字段名
    "CreatedTimeField": "created_time",  //文章表的发布时间字段名，时间戳方式
    "KeywordsField": "keywords",  //文章表的关键词字段名
    "DescriptionField": "description",  //文章表的描述字段名
    "AuthorField": "author",  //文章表的作者字段名
    "ViewsField": "views",  //文章表的浏览量字段名
    "ContentTableName": "fe_new_article_data",  //如果文章内容表和文章表不是同一个表，则在这里填写指定表面，如果相同，则填写相同的名称
    "ContentIdField": "id",  //文章内容表的id字段名
    "ContentField": "content"  //文章内容表或文字表的id字段名
  }
}
```

## 协助完善
欢迎有能力有贡献精神的个人或团体参与到本采集器的开发完善工作中来，共同完善采集功能。请fork一个分支，然后在上面修改，修改完了提交pull request合并请求。

## 版权声明
© Fesion，tpyzlxy@163.com

Released under the [MIT License](https://github.com/fesiong/collector/blob/master/License)