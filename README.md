# kompress

This is a self-learning & proof-og-concept implementation of various lossless compression algorithms, that can be combined into a GZip compatible compressor/decompressor.

By design, it reads the stream of bytes/symbols only once, to comply with the GZip interface.

The API is similar to the gzip golang package : you first construct a writer (resp. a reader) and then write (resp. read) through it to compress (resp. decompress) your data. Writers/Readers can (should) be chained.

It is **essential to close** the Writer when finished, to ensure data is flushed.

For performance, you may want to buffer the initial io.Writer and io.Reader. This is not taken care of by the engines.

You may use these engines to preprocess data before it is gziped, or to replace gzip completely.

At the moment, the following building blocks are :

## MyZip

A typical assembly of a byte-to-symbol block, then a reapeat block, then a lzw block, then a huffann block, then a bit2byte blok. It compresses bytes into bytes.

## DynReader and Writer

The DReader/Writer provides an adaptative huffman compression, compressing symbols into bits (after adding and EOF Symbol to the symbol alphabet). A scheduler defines how frequently the huffman frequency tree is recomputed.

It relies on an *engine* that does the huffan tree management, and hwriter/reader, that implements a fixed tree huffman encoding.

## LZW

This layer will use a dictionnary-based compression, based on the idea of the LZW algorith, to compress from an alphabet to a larger alphabet, buiding a dictionnary of known sequence on the way.

## KDelta 

This layer will not change the alphabet. It tries to predict the next Symbol, based on what it has seen so far, encoding the delta between the prediction and the truth. It does not actually "compress" the message, but improves the statistical properties for a better huffman compression stage if there are some distant redunduncies in the message.

## Repeat

This layer will compress the sequences of identical successive Symbols, using and additionnal "escaped" Symbol. Therefore, the resulting alphabet is one Symbol larger.

## Utilities

**BitBuffer** : A FIFO buffer than can read/write bits, or bytes (seen as 8 bits). Closing triggers a flush, padding with 0 bits. An EOF Symbol must be used to recognized the actual end of file.

**BitFromByteReader/BitToByteWriter** : a conversion layer between bits and bytes.

**LogWriter** : Writes to / reads from  the console, for debugging.

## Note on performance: 

As expected, performance is far from matching the build-in Golang GZip, as can be observed in the provided tests and benchmarks. 

However, some of these blocks maybe used as preprocessing layers, before the built-in GZip is applied, or when processing streams of Symbols that are not bytes, but significantly wider or narrower.

And anyway, it was fun to write and debug !