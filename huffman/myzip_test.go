package huffman

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"testing"
)

func TestMyZipBasic(t *testing.T) {

	source := []byte("Hellooooooo \x00\xFF world !")
	source = append(source, source...)
	source = append(source, source...)

	testReadWriteMyZip(t, source)

	source = []byte("This is standard english text, \x00 et \xff et français, avec des caractères accentués.")
	testReadWriteMyZip(t, source)

	source = []byte("This is standard english text, \x00 et \xff et français, avec des caractères accentués.")
	source = append(source, source...)
	source = append(source, source...)
	testReadWriteMyZip(t, source)

	source = []byte(`
	La nuit était fort noire et la forêt très-sombre.
Hermann à mes côtés me paraissait une ombre.
Nos chevaux galopaient. A la garde de Dieu !
Les nuages du ciel ressemblaient à des marbres.
Les étoiles volaient dans les branches des arbres
Comme un essaim d’oiseaux de feu.

Je suis plein de regrets. Brisé par la souffrance,
L’esprit profond d’Hermann est vide d’espérance.
Je suis plein de regrets. O mes amours, dormez !
Or, tout en traversant ces solitudes vertes,
Hermann me dit : - Je songe aux tombes entr’ouvertes ; -
Et je lui dis : - Je pense aux tombeaux refermés.-

Lui regarde en avant : je regarde en arrière,
Nos chevaux galopaient à travers la clairière ;
Le vent nous apportait de lointains angelus ;
dit : - Je songe à ceux que l’existence afflige,
A ceux qui sont, à ceux qui vivent. - Moi, - lui dis-je,
Je pense à ceux qui ne sont plus !

Les fontaines chantaient. Que disaient les fontaines ?
Les chênes murmuraient. Que murmuraient les chênes ?
Les buissons chuchotaient comme d’anciens amis.
Hermann me dit : - Jamais les vivants ne sommeillent.
En ce moment, des yeux pleurent, d’autres yeux veillent.
Et je lui dis : - Hélas ! d’autres sont endormis !

Hermann reprit alors : - Le malheur, c’est la vie.
Les morts ne souffrent plus. Ils sont heureux ! j’envie
Leur fosse où l’herbe pousse, où s’effeuillent les bois.
Car la nuit les caresse avec ses douces flammes ;
Car le ciel rayonnant calme toutes les âmes
Dans tous les tombeaux à la fois !

Et je lui dis : - Tais-toi ! respect au noir mystère !
Les morts gisent couchés sous nos pieds dans la terre.
Les morts, ce sont les coeurs qui t’aimaient autrefois
C’est ton ange expiré ! c’est ton père et ta mère !
Ne les attristons point par l’ironie amère.
Comme à travers un rêve ils entendent nos voix.	
	`)

	testReadWriteMyZip(t, source)

	source = append(source, source...)
	testReadWriteMyZip(t, source)

}
func testReadWriteMyZip(t *testing.T, source []byte) {

	bb := bytes.NewBuffer(nil)
	w := NewMyZipWriter(bb)

	n, err := w.Write(source)
	if n != len(source) {
		t.Fatal("Could not write till the end !")
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	fmt.Print("MyZip : ", len(source), "\t==> ", len(bb.Bytes()))

	// read back

	r := NewMyZipReader(bb)
	res := make([]byte, len(source))
	n, err = r.Read(res)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(source) {
		t.Fatal("Unexpected length, initial : ", len(source), ", got :", len(res))
	}

	if bytes.Compare(res, source) != 0 {
		fmt.Println("Sent : ", source)
		fmt.Println("Got  : ", res)
		t.Fatal("source and res are not the same ")
	}

	fmt.Println("  \t==> ", len(res))

	// Try one more read
	b := []byte{0}
	_, err = r.Read(b)
	if err != io.EOF {
		t.Fatal("Expected EOF, but got ", err)
	}

	// print gzip for refernec
	testReadWriteGZIP(t, source)

}

// compress with gzip for reference
func testReadWriteGZIP(t *testing.T, source []byte) {

	bb := bytes.NewBuffer(nil)
	g := gzip.NewWriter(bb)

	g.Write(source)
	g.Close()

	fmt.Print("GZip  : ", len(source), "\t==> ", len(bb.Bytes()))

	u, _ := gzip.NewReader(bb)
	res := make([]byte, len(source))

	u.Read(res)
	fmt.Println("  \t==> ", len(res), " \n ")

}
