FROM alpine
RUN mkdir /app & mkdir /conf
ADD ./bin/proxy /bin/
EXPOSE 9090
CMD ["proxy"]
