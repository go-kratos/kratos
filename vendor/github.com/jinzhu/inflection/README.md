Inflection
=========

Inflection pluralizes and singularizes English nouns

## Basic Usage

```go
inflection.Plural("person") => "people"
inflection.Plural("Person") => "People"
inflection.Plural("PERSON") => "PEOPLE"
inflection.Plural("bus")    => "buses"
inflection.Plural("BUS")    => "BUSES"
inflection.Plural("Bus")    => "Buses"

inflection.Singular("people") => "person"
inflection.Singular("People") => "Person"
inflection.Singular("PEOPLE") => "PERSON"
inflection.Singular("buses")  => "bus"
inflection.Singular("BUSES")  => "BUS"
inflection.Singular("Buses")  => "Bus"

inflection.Plural("FancyPerson") => "FancyPeople"
inflection.Singular("FancyPeople") => "FancyPerson"
```

## Register Rules

Standard rules are from Rails's ActiveSupport (https://github.com/rails/rails/blob/master/activesupport/lib/active_support/inflections.rb)

If you want to register more rules, follow:

```
inflection.AddUncountable("fish")
inflection.AddIrregular("person", "people")
inflection.AddPlural("(bu)s$", "${1}ses") # "bus" => "buses" / "BUS" => "BUSES" / "Bus" => "Buses"
inflection.AddSingular("(bus)(es)?$", "${1}") # "buses" => "bus" / "Buses" => "Bus" / "BUSES" => "BUS"
```

## Supporting the project

[![http://patreon.com/jinzhu](http://patreon_public_assets.s3.amazonaws.com/sized/becomeAPatronBanner.png)](http://patreon.com/jinzhu)


## Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).
