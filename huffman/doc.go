// Package huffman provides huffman tree enconding-decoding
// from symboles to bits, or bytes.
// It comes in two flavors : Writers/Reader will compressed based on fixed weights,
// and dWriter / dReader wil compress, updating the huffman tree basedupon actual
// frequencies observed.
// For performance reasons, the tree update is done according to the scheduler function,
// that can be adapted for more flexibility.
package huffman
