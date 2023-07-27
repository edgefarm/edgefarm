import asyncio
import nats
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy, ReplayPolicy


class NatsStream:
    @classmethod
    async def create(
        cls,
        server,
        credsfile_path,
        stream,
        subject,
        durable_name=None,
    ):
        self = NatsStream()

        options = {
            "servers": [server],
        }

        if credsfile_path:
            options["user_credentials"] = credsfile_path

        nc = await nats.connect(**options)

        # Create JetStream context.
        js = nc.jetstream()
        sub = await js.subscribe(subject=subject, durable=durable_name, stream=stream)

        self.sub = sub
        self.nc = nc
        return self

    async def next_msg(self, timeout=2.0):
        """
        Wait for next message.
        Returns: Message object or None if timeout.
        """
        try:
            msg = await self.sub.next_msg(timeout=timeout)
        except asyncio.TimeoutError:
            return None
        return msg

    async def ack(self, msg):
        """
        Acknowledge message.
        """
        await msg.ack()

    async def close(self):
        """
        Close connection.
        """
        await self.nc.close()
