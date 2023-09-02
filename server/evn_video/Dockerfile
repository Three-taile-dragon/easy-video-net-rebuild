FROM alpine
# 维护者
LABEL maintainer="dragon@dragonsss.cn"
WORKDIR /Initial
COPY ./target/evn_video .
RUN apk --no-cache add ffmpeg
RUN  mkdir config
EXPOSE 3459 8884
ENTRYPOINT ["./evn_video"]