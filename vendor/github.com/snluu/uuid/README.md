# uuid

uuid is a lightweight implementation for [Univerally unique identifier](http://en.wikipedia.org/wiki/Universally_unique_identifier).

## Usage

    var id UUID = uuid.Rand()
    fmt.Println(id.Hex())
    fmt.Println(id.Raw())


    id1, err := uuid.FromStr("1870747d-b26c-4507-9518-1ca62bc66e5d")
    id2 := uuid.MustFromStr("1870747db26c450795181ca62bc66e5d")
    fmt.Println(id1 == id2) // true

