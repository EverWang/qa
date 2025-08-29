# 刷刷题后端 Dockerfile
# 使用scratch镜像运行预编译的Linux二进制文件

FROM scratch

# 复制预编译的Linux二进制文件
COPY main-linux /main

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["/main"]