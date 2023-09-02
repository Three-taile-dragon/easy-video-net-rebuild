FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_ws .
RUN  mkdir config
EXPOSE 8081 8883
ENTRYPOINT ["./evn_ws"]