FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_api .
RUN  mkdir config assets
EXPOSE 8080
ENTRYPOINT ["./evn_api"]