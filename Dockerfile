FROM nixos/nix AS base
RUN nix-channel --update

FROM base AS build
# Copy our source and setup our working dir
COPY . /tmp/build
WORKDIR /tmp/build

RUN mkdir -p /output/store
RUN nix --extra-experimental-features "nix-command flakes" profile add --profile /output/profile .#localbox nixpkgs#isolate
RUN cp -R $(nix-store -qR /output/profile) /output/store

FROM base
# Copy over the Nix store and profile from the build stage
COPY --from=build /output/store /nix/store
COPY --from=build /output/profile/ /usr/local/

COPY ./engines /lib/localbox/engines
COPY ./isolate.conf /etc/isolate

ENV SHELL=/bin/sh
ENV ISOLATE_CONFIG_FILE=/etc/isolate

EXPOSE 2000

CMD ["/usr/local/bin/localbox", "-p", "2000"]