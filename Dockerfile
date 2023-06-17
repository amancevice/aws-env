FROM public.ecr.aws/lambda/python
RUN yum install -y golang zip
COPY . /opt/
RUN cd /opt && go build
RUN cd /opt && zip -9r /tmp/package.zip aws-env
ENV AWS_LAMBDA_EXEC_WRAPPER=/opt/aws-env
