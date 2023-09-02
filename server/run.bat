chcp 65001
cd evn_api
docker build -t evn_api:latest .
cd ..
cd evn_article
docker build -t evn_article:latest .
cd ..
cd evn_other
docker build -t evn_other:latest .
cd ..
cd evn_user
docker build -t evn_user:latest .
cd ..
cd evn_video
docker build -t evn_video:latest .
cd ..
cd evn_ws
docker build -t evn_ws:latest .
cd ..
docker-compose up -d