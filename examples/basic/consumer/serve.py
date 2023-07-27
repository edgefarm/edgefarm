import os
import threading
import numpy as np
import asyncio
import json
import base64
from bokeh.io import curdoc
from bokeh.layouts import row, column
from bokeh.models import ColumnDataSource, DataTable, StringEditor, StringFormatter, TableColumn
from bokeh.models import DatetimeTickFormatter, Range1d
from bokeh.plotting import figure
from datetime import datetime, timedelta
from dateutil import tz
from sensor_pb2 import Samples
from stream import NatsStream
import queue

NATS_SERVER = os.environ["NATS_SERVER"]
NATS_EXPORT_SUBJECT = os.environ["NATS_EXPORT_SUBJECT"]
NATS_STREAM_NAME = os.environ["NATS_STREAM_NAME"]
NATS_CONSUMER = os.environ.get("NATS_CONSUMER", None)
NATS_CREDS = os.environ.get("NATS_CREDS", None)


def generate_gradient(start_color, end_color, steps):
    start_color = np.array(start_color)
    end_color = np.array(end_color)
    return [
        tuple(map(int, start_color + (end_color - start_color) * i / (steps - 1)))
        for i in range(steps)
    ]


def hex_to_rgb(value):
    """Return (red, green, blue) for the color given as #rrggbb."""
    value = value.lstrip("#")
    lv = len(value)
    return tuple(int(value[i: i + lv // 3], 16) for i in range(0, lv, lv // 3))


def pb_timestamp_to_local_datetime(pb_timestamp):
    to_zone = tz.tzlocal()
    from_zone = tz.tzutc()
    utc = pb_timestamp.ToDatetime()
    utc = utc.replace(tzinfo=from_zone)
    return utc.astimezone(to_zone)


def localtime_to_utc(localtime):
    return localtime.astimezone(tz.tzutc())


class Display:
    def __init__(self, history_length=timedelta(minutes=1), width=1000):
        self.tablequeue = queue.Queue()

        self.columns = [
            TableColumn(field="timestamps", title="Timestamp",
                        editor=StringEditor(), formatter=StringFormatter()),
            TableColumn(field="node_names", title="Node Name",
                        editor=StringEditor(), formatter=StringFormatter()),
        ]

        self.tableSource = ColumnDataSource(
            data={"timestamps": [], "node_names": []})
        self.data_table = DataTable(source=self.tableSource,
                                    columns=self.columns,
                                    editable=False,
                                    index_position=-1,
                                    index_header="#",
                                    index_width=40,
                                    width=400,
                                    height=200,
                                    selectable=False,
                                    sizing_mode="stretch_width",
                                    )

        self.wave_source = ColumnDataSource(data=dict(x=[], y=[]))
        self.wave_plot = figure(
            title="Waveform",
            sizing_mode="stretch_both",
            y_range=Range1d(-1, 1),
        )
        self.wave_plot.line(
            "x", "y", source=self.wave_source, line_width=2, color="#ff6116"
        )

        self.lock = threading.Lock()
        self.trigger_time = 0
        self.latest_wave = []
        self.last_shown_wave_ts = 0

        self.history_length = history_length
        start, end = self.history_x_range()
        self.history_plot = figure(
            title="Events",
            sizing_mode="stretch_width",
            x_axis_type="datetime",
            x_range=(start, end),
            y_range=(0, 1),
            height=200,
        )
        self.history_data = dict(time=[], amplitude=[])
        self.history_source = ColumnDataSource(
            data=dict(time=[start], amplitude=[0]))
        self.history_plot.vbar(
            x="time",
            top="amplitude",
            width=timedelta(seconds=0.2),
            source=self.history_source,
            color="#ff6116",
        )
        self.history_plot.xaxis.formatter = DatetimeTickFormatter(
            seconds="%H:%M:%S", minutes="%H:%M:%S", minsec="%H:%M:%S", hours="%H:%M:%S"
        )
        self.history_plot.xaxis.major_label_orientation = "horizontal"

        gen_thread = threading.Thread(target=self.read_thread, daemon=True)
        gen_thread.start()
        self.doc = curdoc()
        self.doc.add_next_tick_callback(self.update_wave_plot)
        self.doc.add_next_tick_callback(self.update_table)
        self.doc.add_periodic_callback(self.update_history_plot, 100)
        self.flash_index = 100000
        self.flash_gradient = []
        self.wave_plot_background_fill_color = None
        self.stream = None

    async def read_next_stream_msg(self, stream):
        m = await stream.next_msg(timeout=None)
        js = json.loads(m.data)
        subject = m.subject
        payload = base64.b64decode(js["data_base64"])  # strip dapr header
        data = Samples()
        data.ParseFromString(payload)
        t = pb_timestamp_to_local_datetime(data.triggerTimestamp)
        print("got data len ", len(data.samples), " at ", t)
        x = np.arange(0, len(data.samples), 1)
        y = np.array(data.samples)
        await stream.ack(m)
        return x, y, t, subject

    async def read_loop(self):
        self.stream = await NatsStream.create(
            server=NATS_SERVER,
            credsfile_path=NATS_CREDS,
            stream=NATS_STREAM_NAME,
            subject=NATS_EXPORT_SUBJECT,
            durable_name=NATS_CONSUMER,
        )
        while True:
            x, y, t, subject = await self.read_next_stream_msg(self.stream)
            self.latest_wave = dict(x=x, y=y)
            self.trigger_time = t
            self.history_data["time"].append(self.trigger_time)
            self.history_data["amplitude"].append(np.absolute(y).max())
            new_data = dict(timestamps=[self.trigger_time.strftime(
                "%Y-%d-%m, %H:%M:%S")], node_names=[subject.split(".", 1)[0]])
            self.tablequeue.put(new_data)

    def read_thread(self):
        asyncio.run(self.read_loop())

    def update_table(self):
        while not self.tablequeue.empty():
            item = self.tablequeue.get()
            self.tableSource.stream(item, rollover=200)
        self.doc.add_next_tick_callback(self.update_table)

    def update_wave_plot(self):
        """check at every tick if the wave has been updated, if so, update the plot"""
        if self.trigger_time != self.last_shown_wave_ts:
            with self.lock:
                timestr = self.trigger_time.strftime("%Y-%d-%m, %H:%M:%S")
                self.wave_plot.title.text = f"Waveform captured {timestr}"
                self.wave_plot.y_range.start = -1
                self.wave_plot.y_range.end = 1
                self.history_plot.y_range.start = 0
                self.history_plot.y_range.end = 1

                self.wave_source.data = display.latest_wave
                self.last_shown_wave_ts = self.trigger_time

            if self.wave_plot_background_fill_color is None:
                self.wave_plot_background_fill_color = (
                    self.wave_plot.background_fill_color
                )
                self.flash_gradient = generate_gradient(
                    (255, 255, 255),
                    hex_to_rgb(self.wave_plot_background_fill_color),
                    30,
                )
            self.flash_index = 0

        if self.flash_index < len(self.flash_gradient):
            colcode = "#%02x%02x%02x" % self.flash_gradient[self.flash_index]
            self.wave_plot.background_fill_color = colcode
            self.flash_index += 1
        self.doc.add_next_tick_callback(self.update_wave_plot)

    def update_history_plot(self):
        start, end = self.history_x_range()
        self.history_plot.x_range.start = start
        self.history_plot.x_range.end = end
        hist_times, hist_amplitudes = self.filter_history(end)
        self.history_source.data = dict(
            time=hist_times, amplitude=hist_amplitudes)
        self.history_data["time"] = hist_times
        self.history_data["amplitude"] = hist_amplitudes

    def history_x_range(self):
        now = datetime.now()
        now = now.replace(tzinfo=tz.tzlocal())
        return (now - self.history_length, now)

    def filter_history(self, now):
        """filter out events that are older than history_length"""
        hist_times = []
        hist_amplitudes = []
        with self.lock:
            for t, a in zip(self.history_data["time"], self.history_data["amplitude"]):
                if now - t < self.history_length:
                    hist_times.append(t)
                    hist_amplitudes.append(a)
        return hist_times, hist_amplitudes


display = Display()

layout = column(display.wave_plot,
                row(display.history_plot, display.data_table,
                    sizing_mode="stretch_both"),
                sizing_mode="stretch_both")


curdoc().add_root(layout)
curdoc().title = "Samples"
curdoc().theme = "dark_minimal"
