
to launch servers
-----------------

example:
gnatsd -p 44222 --cluster nats://0.0.0.0:22444 --routes nats://tilient.org:22444,nats://dev.tilient.org:22444 -D

remote example:
ssh -t tilient.org 'tmux new-session -d "gnatsd -p 44222 --cluster nats://0.0.0.0:22444 --routes nats://tilient.org:22444,nats://dev.tilient.org:22444 -D"'

with TLS:
gnatsd -p 44222 --tls --tlscert /etc/ssl/tilient/tilient.org.crt -tlskey /etc/ssl/tilient/tilient.org.key -D

