share
=====

share creates tiny http server that shares selected files with outside world.

Default path to access shared files is $HOST:8080/share/$FILENAME

HTTPS is supported and path/port can be customized with parameters.

Usage
-----

    share [OPTION]... file...

Then access http://$HOST:8080/share/$FILENAME

See --help for additional parameters.

HTTPS
-----

One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
Be sure to add your public ip to the hosts and not just internal one.

    share --cert cert.pem --key key.pem file...

Share should now be accessible at https://$HOST:8080/share/$FILENAME


