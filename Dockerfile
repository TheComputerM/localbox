FROM alpine:latest AS base

RUN wget -qO- https://install.determinate.systems/nix | sh -s -- install linux \
  --extra-conf "sandbox = false" \
  --init none \
  --no-confirm
ENV PATH="${PATH}:/nix/var/nix/profiles/default/bin"

FROM base AS build

COPY . /tmp/build
WORKDIR /tmp/build

RUN mkdir -p /output/store
RUN nix profile add --profile /output/profile .#localbox .#isolate
RUN cp -R $(nix-store -qR /output/profile) /output/store

FROM base

COPY --from=build /output/store /nix/store
COPY --from=build /output/profile/bin/ /usr/local/bin/

COPY ./engines /lib/localbox/engines
COPY ./isolate.conf /etc/isolate

ENV SHELL=/bin/sh
ENV ISOLATE_CONFIG_FILE=/etc/isolate

EXPOSE 2000

CMD ["/usr/local/bin/localbox", "-p", "2000"]