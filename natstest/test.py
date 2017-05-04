import asyncio
from datetime import datetime
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrConnectionClosed, ErrTimeout, ErrNoServers

def run(loop):
  nc = NATS()
  options = {
    "servers": [
      "nats://tilient.org:44222",
      "nats://dev.tilient.org:44222",
    ],
    "io_loop": loop,
  }
  yield from nc.connect(**options)
  yield from nc.publish("count", b"33")
  yield from nc.flush(0.500)
  yield from nc.publish("cmd", b"quit")
  yield from nc.flush(0.500)
  yield from nc.close()

if __name__ == '__main__':
  loop = asyncio.get_event_loop()
  loop.run_until_complete(run(loop))
  loop.close()
