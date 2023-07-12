FROM python:3.10.7-alpine

COPY main.py /usr/local/bin/main.py
COPY requirements.txt /usr/local/bin/requirements.txt
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN apk --update add curl gettext
RUN curl -L -o /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl

RUN chmod a+x /usr/local/bin/kubectl
RUN chmod a+x /usr/local/bin/main.py
RUN chmod a+x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]