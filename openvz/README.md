* Copy ginux-template.tar.gz to /vz/template/cache

* Copy ginux.conf to /etc/vz/dists/

* Copy the contents of scripts/ to /etc/vz/dists/scripts/

* Copy ve-ginux.conf-sample to /etc/vz/conf

To create a new OpenVZ container from this template:

vzctl create ${NUMBER} --config ginux