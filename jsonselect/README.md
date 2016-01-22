jsonselect - Select JSON with class
===================================

Query JSON like you would query some CSS nodes.

jsonselect is a command-line tool to apply JSONSelect filters to
JSON through stdin, line by line.

Usage
-----

Output from `jsonselect --help`:

    Usage of jsonselect:
      -i	Nicely indent any JSON output
      -q	Keep strings quoted instead of unquoting them
      -s	Put things on a single line

From a `jsonfile` that looks like:

    {"event": "My event", "properties": {"os_name": "Windows"}}
    {"event": "My event", "properties": {"os_name": "Linux"}}

Extract the `event` prop, one by line:

    cat jsonfile | jsonselect .event

Extract two lines for each incoming line, one is the `event` property, the other the [JSONPath](http://goessner.net/articles/JsonPath/) equivalent to `.properties.os_name`:

    cat jsonfile | jsonselect .event ".properties .os_name"

Same thing, on a single line, separated by `\t` characters:

    cat jsonfile | jsonselect -s .event ".properties .os_name"

By default, `jsonselect` displays each queried selector one line after
the other.  It is thus possible that one record outputs a variable
number of rows (since some properties can exist only in _some_
records). Using `-s` ensures you have one line per matching record,
but you need to deal with the `\t`.

Nicely indented properties dictionary, prefixed with the `event` as quoted (-q) text:

    cat jsonfile | jsonselect -q -i .event .properties

Merely running `cat jsonfile | jsonselect -i` will display `:root` by default.
