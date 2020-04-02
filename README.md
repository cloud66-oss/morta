<img src="http://cdn2-cloud66-com.s3.amazonaws.com/images/oss-sponsorship.png" width=150/>

# shutdown-sequencer
Given a process PID and a sequence of alternating signals and sleep durations, shutdown-sequencer will perform the sequence against the PID until either the process is dead, or the sequence has completed.

# Installation
Head to the shutdown-sequencer [releases](https://github.com/cloud66/shutdown-sequencer/releases/latest) and download the latest version for your platform.

You can then copy the file to /usr/local/bin and make sure it is renamed to shutdown-sequencer and that it is executable via `chmod +x /usr/local/bin/shutdown-sequencer`. From this point on, you can run `shutdown-sequencer update` to update it automatically.

# Usage
Let's use [unicorn](https://github.com/defunkt/unicorn) as an example, which is an HTTP server for Rack applications. Looking at its [help page for signals](https://github.com/defunkt/unicorn/blob/master/SIGNALS), we can see the following:

Signal | Result
--- | ---
QUIT | graceful shutdown, waits for workers to finish their current request before finishing
INT/TERM | quick shutdown, kills all workers immediately
KILL | terminate process immediately, uncatchable

From this, a reasonable sequence of signals might be the following:
- send SIGQUIT to the master process if it exists
- wait 30 seconds or until master is dead, whichever comes first
- send SIGTERM to the master process if it exists
- wait 10 seconds or until master is dead, whichever comes first
- send SIGKILL to the master process if it exists

Assuming that the master unicorn process has a PID of 1234, you can then run the following to perform the above sequence:
```
shutdown-sequencer -p 1234 -s "quit:30:term:10:kill"
```

You can then use this in your process manager to terminate the process as cleanly as possible. For example, if you're using [systemd](https://www.freedesktop.org/wiki/Software/systemd/), you can add the following to your service definition:
```
ExecStop=-/usr/local/bin/shutdown-sequencer -p $MAINPID -s "quit:30:term:10:kill"
```
