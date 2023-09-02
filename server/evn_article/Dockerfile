FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_article .
RUN  mkdir config
EXPOSE 3460 8885
ENTRYPOINT ["./evn_article"]