FROM 169942020521.dkr.ecr.eu-west-1.amazonaws.com/base/golang:debian11-runtime

CMD ["-bind-address=:6001"]

EXPOSE 6001
