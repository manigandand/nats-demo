#!/bin/bash

STREAM_EVENTS="events"
CONSUMER_CMDR="cmdr"

# create a new stream
# --max-age in microseconds
nats stream add $STREAM_EVENTS 
   --subjects "$STREAM_EVENTS.>"
   --retention work 
   --storage file 
   --discard old 
   --max-msgs=-1 
   --max-bytes=-1 
   --max-age=365
   --max-msg-size=512
   --dupe-window 1h

# create a consumer for the streams
# --max-deliver=2 message is only re tried maximum of 2 times.
# pull based consumer, need explicit acknowledgment
nats consumer add $STREAM_EVENTS $CONSUMER_CMDR 
    --max-deliver=2 
    --pull 
    --wait=5s 
    --ack=explicit 
    --replay=instant 
    --deliver all 
    --filter=""