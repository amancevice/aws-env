FROM public.ecr.aws/lambda/python
RUN yum install -y golang
COPY . /opt/
RUN cd /opt && go build
ENV AWS_LAMBDA_EXEC_WRAPPER=/opt/aws-env
