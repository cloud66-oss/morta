#!/bin/bash

aws s3 cp build s3://downloads.cloud66.com/morta --acl public-read --recursive
