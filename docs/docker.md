docker service set
```
systemctl list-units --type=service

# Check whether startup is set
systemctl list-unit-files | grep enable

# Set boot
systemctl enable docker.service

# Turn off and start up
systemctl disable docker.service

启动时加–restart=always

例如：启动mysql服务,跟随docker一起启动。

docker run -p 3306:3306 --name mysql --restart=always -v /home/mysql/conf:/etc/mysql/conf.d -v /home/mysql/logs:/logs -v /home/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root123 -d mysql:8.1

说明:

-p 3307:3306：将主机的3307端口映射到docker容器的3306端口。
--name mysql：运行服务名字
-v /home/mysql/conf:/etc/mysql/conf.d ：将主机/home/mysql录下的conf/my.cnf 挂载到容器的 /etc/mysql/conf.d
-v /home/mysql/logs:/logs：将主机/home/mysql目录下的 logs 目录挂载到容器的 /logs。
-v /home/mysql/data:/var/lib/mysql ：将主机/home/mysql目录下的data目录挂载到容器的 /var/lib/mysql 
-e MYSQL_ROOT_PASSWORD=root123：初始化 root 用户的密码。
-d mysql:5.6 : 后台程序运行mysql816

已经运行的容器可更新为：
docker update --restart=always 容器id或name
```