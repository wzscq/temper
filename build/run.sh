docker run -d --name temperpredictservice -p8004:80 -v /root/temperpredictservice/conf:/services/temperpredictservice/conf  wangzhsh/temperpredictservice:0.0.1

1、 install python3.9
      按照依赖
      yum install zlib-devel bzip2-devel openssl-devel ncurses-devel readline-devel tk-devel gcc make libffi-devel
      下载按照包
      wget https://www.python.org/ftp/python/3.9.17/Python-3.9.17.tgz
      tar -xvf Python-3.9.17.tgz
      cd Python-3.9.17
      ./configure
      make
      make install