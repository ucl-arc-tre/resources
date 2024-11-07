# Adapted from works in https://github.com/mathworks-ref-arch/matlab-dockerfile
# Copyright 2023-2024 The MathWorks, Inc.
ARG MATLAB_RELEASE=r2023a
FROM mathworks/matlab:${MATLAB_RELEASE}

USER matlab

# When running the container a license file can be mounted,
# or a license server can be provided as an environment variable.
# For more information, see https://hub.docker.com/r/mathworks/matlab

WORKDIR /home/matlab
# ENTRYPOINT and CMD are inherited from the base image
