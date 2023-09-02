FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_other .
RUN  mkdir config
EXPOSE 3458 8882
ENTRYPOINT ["./evn_other"]