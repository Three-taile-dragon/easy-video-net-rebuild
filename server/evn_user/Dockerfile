FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_user .
RUN  mkdir config
EXPOSE 3457 8881
ENTRYPOINT ["./evn_user"]