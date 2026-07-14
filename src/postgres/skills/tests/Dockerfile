# Docker 测试环境 —— 模拟一个完全干净的新用户机器
# 用法: docker build -t pg-skill-test . && docker run --rm pg-skill-test

FROM python:3.11-slim

# 工作目录
WORKDIR /app/test

# 只装测试脚本需要的最小依赖（不含 tencentcloud-sdk）
COPY test_new_user.py /app/test/

# 把打包好的 zip 拷进去（在 docker build 时复制）
COPY dist/*.zip /app/test/

# 关键：不设置任何环境变量，模拟"零配置新用户"
ENV TEST_MODE=container

ENTRYPOINT ["python3", "test_new_user.py"]
