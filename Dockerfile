FROM public.ecr.aws/lambda/python
RUN yum install -y golang
COPY . /opt
RUN cd /opt && go build
RUN mv /opt/cmd/index.py /var/task/index.py
