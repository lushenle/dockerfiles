FROM ajoergensen/openssh-server:latest
LABEL maintainer="Shenle Lu <lushenle@gmail.com>" run="openssh-server" app="tmux"
RUN apk update \
    && apk --no-cache add tmux zip vim ca-certificates \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/* \
    && rm -r /usr/share/man
ADD .tmux.conf ~/

