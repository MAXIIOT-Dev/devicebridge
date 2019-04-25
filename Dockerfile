FROM busybox 
 
WORKDIR /usr/local/devicebridge/

# RUN mkdir config
RUN mkdir ui
COPY ./build/devicebridge .
 
ENTRYPOINT ["./devicebridge"]