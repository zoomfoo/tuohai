FROM ubuntu:14.04
MAINTAINER wangyang tosingular@gmail.com


RUN mkdir -p /home/tuohai/apps && cd /home/tuohai/apps &&mkdir file_api im_api open_api
ADD URL/file_api /home/tuohai/apps

CMD ["sh run.sh","all"]



