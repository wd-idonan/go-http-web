# go-http-web
基于http重写的框架

# 注意
连接mysql 8.0以上版本会认证失败
this authentication plugin is not supported

# 修改
mysql -u root -p
use mysql;
select user,host from user;
+------------------+-----------+
| user             | host      |
+------------------+-----------+
| root             | %         |
| mysql.infoschema | localhost |
| mysql.session    | localhost |
| mysql.sys        | localhost |
+------------------+-----------+
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '$password'