FROM centos:7

WORKDIR /github/workspace

RUN yum install -y rpmdevtools yum-utils spectool && \
    yum clean all && \
    rm -r -f /var/cache/*

CMD spectool --get-files --all SPECS/*.spec && \
    yum-builddep --assumeyes SPECS/*.spec && \
    rpmbuild \
        --define '_topdir /github/workspace' \
        --define "_version $(echo $GITHUB_REF | cut -d'/' -f3)" \
        --define "_dist .el7" \
        -ba SPECS/*.spec
