FROM byuoitav/amd64-alpine
MAINTAINER Daniel Randall <danny_randall@byu.edu>

COPY pi-designation.monsters pi-designation.monsters
COPY version.txt version.txt

ENTRYPOINT ./pi-designation.monsters

