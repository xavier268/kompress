# kompress

Compression engines (complementing / improving go-gzip package)

The API is similar to the gzip golang package : you first construct a writer (resp. a reader) and then write (resp. read) through it to compress (resp. decompress) your data.

You may use these engines to preprocess data before it is gziped, or to replace gzip completely.

At the moment, the following engines are provided :

## Klog

Not actually a compression. 

Does not change the bytes. Just dump on screen the bytes as they flow through.

## Krlen

Local compression. 

It replaces sequences of identical bytes with specially encoded shorter sequences, using an escape byte. The escape byte is dynamically selected 
from the least frequent bytes.

## Kdelta

Global transformation. 

Not really a compression, it tries to predict the next bytes, based on the history,
and will encode the delta. This is a way to detect possibly long sequences far away in the file, and encode them in a way that will require only local optimization (Krlen and/or Kbit) as a next step.

## Kbit (TO DO)

Local compression.

It encodes the bytes using a variable length system.
To select how to encode, it dynamically adjust based on bytes frequency.


