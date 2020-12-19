############################
# STEP 2 build a small image
############################
FROM alpine

LABEL maintainer "automsen Community <sf@v2ray.com>"

RUN echo "http://mirrors.aliyun.com/alpine/v3.12/main/" > /etc/apk/repositories
RUN apk update && apk add --no-cache ca-certificates bash wget curl && \
    mkdir -p /usr/bin/v2ray
ADD * /usr/bin/v2ray

#ENTRYPOINT ["/usr/bin/v2ray/v2ray"]
ENV PATH /usr/bin/v2ray:$PATH
CMD ["v2ray", "-config=/etc/v2ray/config.json"]
