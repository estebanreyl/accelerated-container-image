FROM ubuntu:latest AS base
# Required Build/Run Tools Dependencies for Overlaybd tools
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates

RUN apt update && \
    apt install -y libcurl4-openssl-dev libext2fs-dev libaio-dev

# --- OVERLAYBD TOOLS ---
FROM base As overlaybd-build
RUN apt update && \
    apt install -y libssl-dev libnl-3-dev libnl-genl-3-dev libgflags-dev libzstd-dev && \
    apt install -y zlib1g-dev binutils && \
    apt install -y make git wget sudo tar gcc cmake && \
    apt install -y golang

RUN git clone https://github.com/containerd/overlaybd.git && \
    cd overlaybd && \
    git submodule update --init && \
    git checkout fc255f39800ba01e80ae414514ac953ddced4842 && \
    mkdir build && \
    cd build && \
    cmake .. && \
    make -j && \
    make install

# --- BUILD LOCAL CONVERTER ---
FROM overlaybd-build AS convert-build
WORKDIR /home/limiteduser/
RUN git clone https://github.com/containerd/accelerated-container-image.git
WORKDIR /home/limiteduser/accelerated-container-image
RUN make

# --- FINAL ---
FROM base
WORKDIR /home/limiteduser/

# Copy Conversion Tools
COPY --from=overlaybd-build /opt/overlaybd/bin /opt/overlaybd/bin
COPY --from=overlaybd-build /opt/overlaybd/baselayers /opt/overlaybd/baselayers

# # This is necessary for overlaybd_apply to work
COPY --from=overlaybd-build /etc/overlaybd/overlaybd.json /etc/overlaybd/overlaybd.json

COPY --from=convert-build /home/limiteduser/accelerated-container-image/bin/convertor ./bin/convertor
CMD ["./bin/convertor"]