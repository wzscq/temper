FROM alpine:3.15
RUN set -eux &&sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add nginx &&\ 
    apk add ffmpeg && \
    apk update && \
    apk add --no-cache musl-dev gcc python3 python3-dev py3-pip && \
    python3 -m venv /venv && \
    . /venv/bin/activate  && \
    pip3 install keras numpy sklearn pandas  && \
    deactivate

COPY nginx.conf /etc/nginx/nginx.conf
#ADD ./package/web /web
ADD ./package/service/audioscoring /services/audioscoringservice/audioscoring
ADD audioscore.py /services/audioscoringservice/audioscore.py
copy entrypoint.sh /entrypoint.sh
ENTRYPOINT ["sh","entrypoint.sh"]