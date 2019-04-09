FROM busybox 
 
WORKDIR /usr/local/vbaseBridge/

RUN mkdir config
RUN mkdir static
COPY ./build/vbasebridge .
COPY ./ui ./ui

# COPY ./config.toml /usr/local/NursingHome/
ENTRYPOINT ["./vbasebridge"]