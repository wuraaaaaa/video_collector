# 某视频CMS资源站数据入库采集

数据库采用MySQL 按需求创建.env文件，参照.env.example
采用协程池调度采集(需修改mysql配置防止速度过快丢失)