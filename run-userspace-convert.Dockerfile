FROM mcr.microsoft.com/cbl-mariner/base/core:2.0 AS base
# Required Build/Run Tools Dependencies for Overlaybd tools
RUN yum install e2fsprogs-devel -y && \
    yum install libaio-devel -y && \
    yum install ca-certificates -y && \
    yum install shadow-utils -y

# --- OVERLAYBD TOOLS ---
FROM base As overlaybd-build
RUN yum install -y libaio-devel libcurl-devel openssl-devel libnl3-devel e2fsprogs-devel glibc-devel libzstd-devel binutils ca-certificates-microsoft build-essential && \
    yum install -y rpm-build make git wget sudo tar gcc gcc-c++ cmake && \
    yum install golang -y

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

# This is necessary for overlaybd_apply to work
COPY --from=overlaybd-build /etc/overlaybd/overlaybd.json /etc/overlaybd/overlaybd.json

COPY --from=convert-build /home/limiteduser/accelerated-container-image/bin/convertor ./bin/convertor
CMD ["./bin/convertor"]