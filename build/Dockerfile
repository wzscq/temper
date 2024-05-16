FROM centos:7
RUN yum install -y zlib-devel bzip2-devel openssl-devel ncurses-devel readline-devel tk-devel gcc make libffi-devel wget epel-release nginx && \
    wget https://www.python.org/ftp/python/3.9.17/Python-3.9.17.tgz && \
    tar -zxvf Python-3.9.17.tgz && \
    cd Python-3.9.17 && ./configure && make && make install
RUN pip3 uninstall -y urllib3 && \
    pip3 install urllib3==1.23 keras tensorflow  scikit-learn pandas

COPY nginx.conf /etc/nginx/nginx.conf
#ADD ./package/web /web
ADD ./package/service/temper /root/services/temperpredictservice/temper
ADD predict.py /root/services/temperpredictservice/predict.py
copy entrypoint.sh /root/entrypoint.sh
ENTRYPOINT ["sh","/root/entrypoint.sh"]